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

	var data string

	for index := 0; index < len(function); index++ {
		currentData, err := contract.getHexValue(function[index], args[index])

		if err != nil {
			return nil, err
		}

		data += currentData
	}

	transaction.Data = types.ComplexString(sha3Function[0:10] + data)

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
