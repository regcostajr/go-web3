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

package web3

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cellcycle/go-web3/complex/types"
	"github.com/cellcycle/go-web3/dto"
	"golang.org/x/crypto/sha3"
	"regexp"
	"strconv"
	"strings"

	"math/big"
)

// Contract ...
type Contract struct {
	super     *Eth
	abi       string
	functions map[string][]string
}

// NewContract - Contract abstraction
func (eth *Eth) NewContract(abi string) (*Contract, error) {

	contract := new(Contract)
	var mockInterface interface{}

	err := json.Unmarshal([]byte(abi), &mockInterface)

	if err != nil {
		return nil, err
	}

	jsonInterface := mockInterface.([]interface{})
	contract.functions = make(map[string][]string)
	for index := 0; index < len(jsonInterface); index++ {
		function := jsonInterface[index].(map[string]interface{})

		if function["type"] == "constructor" || function["type"] == "fallback" {
			function["name"] = function["type"]
		}

		functionName := function["name"].(string)
		contract.functions[functionName] = make([]string, 0)

		if function["inputs"] == nil {
			continue
		}

		inputs := function["inputs"].([]interface{})
		for paramIndex := 0; paramIndex < len(inputs); paramIndex++ {
			params := inputs[paramIndex].(map[string]interface{})
			contract.functions[functionName] = append(contract.functions[functionName], params["type"].(string))
		}

	}

	contract.abi = abi
	contract.super = eth

	return contract, nil
}

// prepareTransaction ...
func (contract *Contract) prepareTransaction(transaction *dto.TransactionParameters, functionName string, args []interface{}) (*dto.TransactionParameters, error) {

	function, ok := contract.functions[functionName]
	if !ok {
		return nil, errors.New("Function not found")
	}

	fullFunction := functionName + "("

	comma := ""
	for arg := range function {
		fullFunction += comma + function[arg]
		comma = ","
	}

	fullFunction += ")"

	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(fullFunction))
	sha3Function := fmt.Sprintf("0x%x", hash.Sum(nil))

	var static []string
	var dynamic []string

	offsetCount := contract.calculateOffset(function)
	offset := offsetCount * 32

	for index := 0; index < len(function); index++ {
		currentData, err := contract.encode(function[index], args[index])

		if err != nil {
			return nil, err
		}

		if contract.isDynamic(function[index]) {
			hexOffset, _ := contract.encodeUint(offset, "")
			static = append(static, hexOffset[0])
			dynamic = append(dynamic, currentData...)
			offset = offset + 32
		}

		static = append(static, currentData...)
	}

	static = append(static, dynamic...)

	transaction.Data = types.ComplexString(sha3Function[0:10] + strings.Join(static, ""))

	return transaction, nil

}

func (contract *Contract) Call(transaction *dto.TransactionParameters, functionName string, args ...interface{}) (*dto.RequestResult, error) {

	transaction, err := contract.prepareTransaction(transaction, functionName, args)

	if err != nil {
		return nil, err
	}

	return contract.super.Call(transaction)

}

func (contract *Contract) Send(transaction *dto.TransactionParameters, functionName string, args ...interface{}) (string, error) {

	transaction, err := contract.prepareTransaction(transaction, functionName, args)

	if err != nil {
		return "", err
	}

	return contract.super.SendTransaction(transaction)

}

func (contract *Contract) Deploy(transaction *dto.TransactionParameters, bytecode string, args ...interface{}) (string, error) {

	constructor := contract.functions["constructor"]

	for index := 0; index < len(constructor); index++ {
		tmpBytes, err := contract.getHexValue(constructor[index], args[index])

		if err != nil {
			return "", err
		}

		bytecode += tmpBytes
	}

	transaction.Data = types.ComplexString(bytecode)

	return contract.super.SendTransaction(transaction)

}

func (contract *Contract) calculateOffset(args []string) int {
	offset := 0

	for index := 0; index < len(args); index++ {
		regex := regexp.MustCompile(`^(?:[a-z]+)(\d+)?(?:(?:\[)(\d+)(?:\]))?`)
		match := regex.FindStringSubmatch(args[index])

		itemSize := match[1]
		arraySize := match[2]

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

func (contract *Contract) isDynamic(inputType string) bool {

	regex := regexp.MustCompile(`(\[\])`)
	match := regex.FindStringSubmatch(inputType)

	// non fixed size array
	if len(match) > 0 && match[1] != "" {
		return true
	}

	regex = regexp.MustCompile(`^(address|uint|int|ufixed|fixed|bool)`)
	match = regex.FindStringSubmatch(inputType)

	if len(match) > 0 && match[1] != "" {
		return false
	}

	regex = regexp.MustCompile(`^(string|bytes)(\d+)?`)
	match = regex.FindStringSubmatch(inputType)

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
	}

	return methodMap[function]
}

func (contract *Contract) encode(inputType string, value interface{}) ([]string, error) {
	regex := regexp.MustCompile(`^([a-z]+)(\d+)?(\[\d+\])?`)
	match := regex.FindStringSubmatch(inputType)

	basicType := match[1]
	itemSize := match[2]
	array := match[3]

	// array
	if array != "" {
		arrayValues := value.([]interface{})
		var s []string
		for _, v := range arrayValues {
			encoded, err := contract.encodeMap(basicType).(func(interface{}, string) ([]string, error))(v, itemSize)
			if err != nil {
				return nil, err
			}
			s = append(s, encoded...)
		}
	}

	return contract.encodeMap(basicType).(func(interface{}, string) ([]string, error))(value, itemSize)
}

func (contract *Contract) encodeString(value interface{}, _ string) ([]string, error) {
	var s []string

	size := fmt.Sprintf("%064s", fmt.Sprintf("%x", len(value.(string))))
	s = append(s, size)

	hex := fmt.Sprintf("%x", value.(string))
	hex += strings.Repeat("0", 64-len(hex))
	s = append(s, hex)

	return s, nil
}

func (contract *Contract) encodeUint(value interface{}, size string) ([]string, error) {
	bigValue := value.(*big.Int)
	if bigValue.Cmp(big.NewInt(0)) == -1 {
		return nil, errors.New(fmt.Sprintf("Int type lower than 0: %s", bigValue.String()))
	}
	return contract.encodeInt(value, size)
}

func (contract *Contract) encodeInt(value interface{}, size string) ([]string, error) {
	bigValue := value.(*big.Int)
	if size != "" {
		intSize, err := strconv.Atoi(size)
		if err != nil {
			return nil, errors.New("Invalid size for input type, please check the ABI for typos")
		}

		if bigValue.BitLen() > intSize {
			return nil, errors.New(fmt.Sprintf("Input type size does not match with the ABI information: %s, ABI: %d", bigValue.String(), intSize))
		}
	}
	return []string{fmt.Sprintf("%064s", fmt.Sprintf("%x", bigValue.String()))}, nil
}

func (contract *Contract) encodeAddress(value interface{}, _ string) ([]string, error) {
	// removes 0x
	return []string{fmt.Sprintf("%064s", value.(string)[2:])}, nil
}

func (contract *Contract) getHexValue(inputType string, value interface{}) (string, error) {

	var data string

	if strings.HasPrefix(inputType, "int") ||
		strings.HasPrefix(inputType, "uint") ||
		strings.HasPrefix(inputType, "fixed") ||
		strings.HasPrefix(inputType, "ufixed") {

		bigVal := value.(*big.Int)

		// Checking that the string actually is the correct inputType
		if strings.Contains(inputType, "128") {
			// 128 bit
			if bigVal.BitLen() > 128 {
				return "", errors.New(fmt.Sprintf("Input type %s not met", inputType))
			}
		} else if strings.Contains(inputType, "256") {
			// 256 bit
			if bigVal.BitLen() > 256 {
				return "", errors.New(fmt.Sprintf("Input type %s not met", inputType))
			}
		}

		data += fmt.Sprintf("%064s", fmt.Sprintf("%x", bigVal.String()))
	}

	if strings.Compare("address", inputType) == 0 {
		data += fmt.Sprintf("%064s", value.(string)[2:])
	}

	if strings.Compare("string", inputType) == 0 {
		data += fmt.Sprintf("%064s", fmt.Sprintf("%x", value.(string)))
	}

	return data, nil

}
