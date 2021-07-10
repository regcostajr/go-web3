/********************************************************************************
   This file is part of go-web3.
   go-web3 is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-web3 is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-web3.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

/**
 * @file contract.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2018
 */

package eth

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"math/big"

	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/utils"
	"golang.org/x/crypto/sha3"
)

const CONTRACT_CONSTRUCTOR = "constructor"

// Contract ...
type Contract struct {
	super       *Eth
	abi         []ABIItem
	Functions   func(string, int) *ABIItem
	Events      func(string, int) *ABIItem
	Constructor func(int) *ABIItem
}

type ABIItem struct {
	Inputs          []ABIItemField `json:"inputs"`
	Outputs         []ABIItemField `json:"outputs"`
	StateMutability string         `json:"stateMutability"`
	Type            string         `json:"type"`
	Anonymouns      bool           `json:"anonymous"`
	Name            string         `json:"name"`
	Call            func(transaction *dto.TransactionParameters, args ...interface{}) (*dto.RequestResult, error)
	Send            func(transaction *dto.TransactionParameters, args ...interface{}) (string, error)
}

type ABIItemField struct {
	Indexed      bool   `json:"indexed"`
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}

// NewContract - Contract abstraction
func (eth *Eth) NewContract(abi string) (*Contract, error) {

	contract := new(Contract)
	var abiItems []ABIItem

	err := json.Unmarshal([]byte(abi), &abiItems)

	if err != nil {
		return nil, err
	}

	contract.Functions = func(name string, index int) *ABIItem {
		function := contract.find(name, index, "function")
		function.Call = func(transaction *dto.TransactionParameters, args ...interface{}) (*dto.RequestResult, error) {
			return contract.Call(transaction, function, args...)
		}
		function.Send = func(transaction *dto.TransactionParameters, args ...interface{}) (string, error) {
			return contract.Send(transaction, function, args...)
		}
		return function
	}
	contract.Events = func(name string, index int) *ABIItem {
		return contract.find(name, index, "event")
	}
	contract.Constructor = func(index int) *ABIItem {
		return contract.find(CONTRACT_CONSTRUCTOR, index, CONTRACT_CONSTRUCTOR)
	}

	contract.abi = abiItems
	contract.super = eth

	return contract, nil
}

func (contract *Contract) find(name string, index int, sType string) *ABIItem {
	var s []ABIItem
	for _, v := range contract.abi {
		if (v.Name == name || v.Type == name) && v.Type == sType {
			s = append(s, v)
		}
	}

	if len(s) > 0 {
		return &s[index]
	} else {
		return nil
	}
}

func (contract *Contract) GetSHA3FunctionString(function *ABIItem) string {
	var sb strings.Builder

	sb.WriteString(function.Name)
	sb.WriteString("(")

	comma := ""
	for _, v := range function.Inputs {
		sb.WriteString(comma + v.Type)
		comma = ","
	}

	sb.WriteString(")")
	fullFunction := sb.String()

	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(fullFunction))
	return fmt.Sprintf("0x%x", hash.Sum(nil))
}

// PrepareTransaction ...
func (contract *Contract) PrepareTransaction(transaction *dto.TransactionParameters, function *ABIItem, bytecode []byte, args ...interface{}) (*dto.TransactionParameters, error) {
	// must be cleaned in case of this is called twice for the same transaction
	transaction.Data = ""

	var functionID string
	if function.Type != CONTRACT_CONSTRUCTOR {
		functionID = contract.GetSHA3FunctionString(function)[0:10]
	}

	var static []string
	var dynamic []string

	offsetCount := contract.calculateOffset(function.Inputs)
	offset := offsetCount * 32
	internalArrayOffset := 0

	for index := 0; index < len(function.Inputs); index++ {
		currentData, err := contract.encode(function.Inputs[index], args[index], &internalArrayOffset)

		if err != nil {
			return nil, err
		}

		if contract.isDynamic(function.Inputs[index]) {
			hexOffset, _ := contract.encodeUint(big.NewInt(int64(offset)), "")
			static = append(static, hexOffset[0])
			dynamic = append(dynamic, currentData...)
			offset = offset + (32 * len(currentData))
		} else {
			static = append(static, currentData...)
		}
	}

	static = append(static, dynamic...)

	strBytecode := string(bytecode)
	strBytecode = strings.TrimSuffix(strBytecode, "\n")
	strBytecode = strings.TrimPrefix(strBytecode, "0x")
	functionID = strings.TrimPrefix(functionID, "0x")
	transaction.Data = fmt.Sprintf("0x%s%s%s", strBytecode, functionID, strings.Join(static, ""))

	return transaction, nil
}

func (contract *Contract) Call(transaction *dto.TransactionParameters, function *ABIItem, args ...interface{}) (*dto.RequestResult, error) {

	transaction, err := contract.PrepareTransaction(transaction, function, nil, args...)

	if err != nil {
		return nil, err
	}

	return contract.super.Call(transaction)

}

func (contract *Contract) Send(transaction *dto.TransactionParameters, function *ABIItem, args ...interface{}) (string, error) {

	transaction, err := contract.PrepareTransaction(transaction, function, nil, args...)

	if err != nil {
		return "", err
	}

	return contract.super.SendTransaction(transaction)

}

func (contract *Contract) Deploy(transaction *dto.TransactionParameters, bytecode []byte, args ...interface{}) (string, error) {
	transaction, err := contract.PrepareTransaction(transaction, contract.Constructor(0), bytecode, args...)

	if err != nil {
		return "", err
	}

	return contract.super.SendTransaction(transaction)
}

func (contract *Contract) calculateOffset(args []ABIItemField) int {
	offset := 0

	for index := 0; index < len(args); index++ {
		regex := regexp.MustCompile(`^([a-z]+)(\d+)?(?:(?:\[)(\d+)(?:\]))?`)
		match := regex.FindStringSubmatch(args[index].Type)

		itemSize := match[2]
		arraySize := match[3]

		// fixed array size and fixed var size
		if itemSize != "" && arraySize != "" {
			i, _ := strconv.Atoi(arraySize)
			offset += i
			continue
		}

		offset++
	}

	return offset
}

func (contract *Contract) isDynamic(input ABIItemField) bool {

	regex := regexp.MustCompile(`(\[\])`)
	match := regex.FindStringSubmatch(input.Type)

	// non fixed size array
	if len(match) > 0 && match[1] != "" {
		return true
	}

	regex = regexp.MustCompile(`^(address|uint|int|bool)`)
	match = regex.FindStringSubmatch(input.Type)

	if len(match) > 0 && match[1] != "" {
		return false
	}

	regex = regexp.MustCompile(`^(string|bytes)(\d+)?`)
	match = regex.FindStringSubmatch(input.Type)

	if len(match) > 1 && match[2] != "" {
		return false
	}

	return true

}

func (contract *Contract) encodeMap(function string) interface{} {
	methodMap := map[string]interface{}{
		"string":  contract.encodeString,
		"int":     contract.encodeInt,
		"uint":    contract.encodeUint,
		"address": contract.encodeAddress,
		"bytes":   contract.encodeBytes,
	}

	return methodMap[function]
}

func (contract *Contract) encode(input ABIItemField, value interface{}, internalArrayOffset *int) ([]string, error) {
	regex := regexp.MustCompile(`^([a-z]+)(\d+)?(\[(\d+)?\])?`)
	match := regex.FindStringSubmatch(input.Type)

	basicType := match[1]
	itemSize := match[2]
	array := match[3]

	// array
	if array != "" {
		arrayValue := reflect.ValueOf(value)
		arrayValues := make([]interface{}, arrayValue.Len())
		for i := 0; i < arrayValue.Len(); i++ {
			arrayValues[i] = arrayValue.Index(i).Interface()
		}

		s := make([]string, 0)
		if array == "[]" {
			arraySize, _ := contract.encodeUint(big.NewInt(int64(len(arrayValues))), "")
			s = append(s, arraySize[0])
		}

		var static []string
		var dynamic []string

		for k, v := range arrayValues {
			encoded, err := contract.encodeMap(basicType).(func(interface{}, string) ([]string, error))(v, itemSize)
			if err != nil {
				return nil, err
			}

			if contract.isDynamic(ABIItemField{Type: basicType + itemSize}) {
				correctIndex := k + 1
				if internalArrayOffset != nil {
					correctIndex = correctIndex + *internalArrayOffset
					*internalArrayOffset = correctIndex
				}

				hexOffset, _ := contract.encodeUint(big.NewInt(int64(correctIndex*32)), "")
				static = append(static, hexOffset[0])
			}
			dynamic = append(dynamic, encoded...)
		}

		s = append(s, static...)
		s = append(s, dynamic...)
		return s, nil
	}

	return contract.encodeMap(basicType).(func(interface{}, string) ([]string, error))(value, itemSize)
}

func (contract *Contract) encodeString(value interface{}, _ string) ([]string, error) {
	var s []string

	size := fmt.Sprintf("%064s", fmt.Sprintf("%x", len(value.(string))))
	s = append(s, size)

	hexString := fmt.Sprintf("%x", value.(string))
	hexString += strings.Repeat("0", 64-len(hexString))
	s = append(s, hexString)

	return s, nil
}

func (contract *Contract) DecodeString(hexString string) (string, error) {
	data, err := hex.DecodeString(hexString)
	if err != nil {
		return "", err
	}

	return utils.CleanString(string(data)), nil
}

func (contract *Contract) getBigIntFromFlexibleParameter(value interface{}) *big.Int {
	iKind := reflect.ValueOf(value).Kind()

	var bigValue *big.Int
	if strings.HasPrefix(iKind.String(), "int") {
		bigValue = big.NewInt(value.(int64))
	} else {
		bigValue = value.(*big.Int)
	}

	return bigValue
}

func (contract *Contract) encodeUint(value interface{}, size string) ([]string, error) {
	bigValue := contract.getBigIntFromFlexibleParameter(value)

	if bigValue.Cmp(big.NewInt(0)) == -1 {
		return nil, errors.New(fmt.Sprintf("Int type lower than 0: %s", bigValue.String()))
	}

	return contract.encodeInt(value, size)
}

func (contract *Contract) DecodeUint(hexString string) *big.Int {
	return contract.DecodeInt(hexString)
}

func (contract *Contract) encodeInt(value interface{}, size string) ([]string, error) {
	bigValue := contract.getBigIntFromFlexibleParameter(value)

	if size != "" {
		intSize, err := strconv.Atoi(size)
		if err != nil {
			return nil, errors.New("Invalid size for input type, please check the ABI for typos")
		}

		if bigValue.BitLen() > intSize {
			return nil, errors.New(fmt.Sprintf("Input type size does not match with the ABI information: %s, ABI: %d", bigValue.String(), intSize))
		}
	}
	return []string{fmt.Sprintf("%064s", fmt.Sprintf("%x", bigValue))}, nil
}

func (contract *Contract) DecodeInt(hexString string) *big.Int {
	big := new(big.Int)
	big.SetString(hexString, 16)
	return big
}

func (contract *Contract) encodeAddress(value interface{}, _ string) ([]string, error) {
	// removes 0x
	return []string{fmt.Sprintf("%064s", value.(string)[2:])}, nil
}

func (contract *Contract) DecodeAddress(address string) string {
	return utils.CleanString(fmt.Sprintf("0x%s", address[len(address)-40:]))
}

func (contract *Contract) encodeBytes(value interface{}, size string) ([]string, error) {
	var s []string

	byteArray := value.([]byte)

	if size == "" {
		bSize := fmt.Sprintf("%064s", fmt.Sprintf("%x", len(byteArray)))
		s = append(s, bSize)
	}

	hexString := hex.EncodeToString(byteArray)
	hexString += strings.Repeat("0", 64-len(hexString))
	s = append(s, hexString)

	return s, nil
}

func (contract *Contract) decodeBytes(hexString string) ([]byte, error) {
	strBytes, err := contract.DecodeString(hexString)

	if err != nil {
		return nil, err
	}

	return []byte(strBytes), nil
}
