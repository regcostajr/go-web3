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
 * @file transaction.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package dto

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/cellcycle/go-web3/utils"
)

// TransactionParameters GO transaction to make more easy control the parameters
type TransactionParameters struct {
	From     string
	To       string
	Nonce    *big.Int
	Gas      *big.Int
	GasPrice *big.Int
	Value    *big.Int
	Data     string
}

// RequestTransactionParameters JSON
type RequestTransactionParameters struct {
	From     string `json:"from"`
	To       string `json:"to,omitempty"`
	Nonce    string `json:"nonce,omitempty"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Value    string `json:"value,omitempty"`
	Data     string `json:"data,omitempty"`
}

// Transform the GO transactions parameters to json style
func (params *TransactionParameters) Transform() *RequestTransactionParameters {
	request := new(RequestTransactionParameters)
	request.From = params.From
	if params.To != "" {
		request.To = params.To
	}
	if params.Nonce != nil {
		request.Nonce = "0x" + params.Nonce.Text(16)
	}
	if params.Gas != nil {
		request.Gas = "0x" + params.Gas.Text(16)
	}
	if params.GasPrice != nil {
		request.GasPrice = "0x" + params.GasPrice.Text(16)
	}
	if params.Value != nil {
		request.Value = "0x" + params.Value.Text(16)
	}
	if params.Data != "" {
		if strings.HasPrefix(params.Data, "0x") {
			request.Data = params.Data
		} else {
			request.Data = "0x" + hex.EncodeToString([]byte(params.Data))
		}
	}
	return request
}

type TransactionResponse struct {
	Hash             string   `json:"hash"`
	Nonce            *big.Int `json:"nonce"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      *big.Int `json:"blockNumber"`
	TransactionIndex *big.Int `json:"transactionIndex"`
	From             string   `json:"from"`
	To               string   `json:"to"`
	Input            string   `json:"input"`
	Value            *big.Int `json:"value"`
	GasPrice         *big.Int `json:"gasPrice,omitempty"`
	Gas              *big.Int `json:"gas,omitempty"`
	Data             *string  `json:"data,omitempty"`
}

type TransactionReceipt struct {
	TransactionHash   string           `json:"transactionHash"`
	TransactionIndex  *big.Int         `json:"transactionIndex"`
	BlockHash         string           `json:"blockHash"`
	BlockNumber       *big.Int         `json:"blockNumber"`
	From              string           `json:"from"`
	To                string           `json:"to"`
	CumulativeGasUsed *big.Int         `json:"cumulativeGasUsed"`
	GasUsed           *big.Int         `json:"gasUsed"`
	ContractAddress   string           `json:"contractAddress"`
	Logs              []TransactionLog `json:"logs"`
	LogsBloom         string           `json:"logsBloom"`
	Status            bool             `json:"status"`
	Type              *big.Int         `json:"type"`
}

type TransactionLog struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      *big.Int `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex *big.Int `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         *big.Int `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

func (t *TransactionResponse) UnmarshalJSON(data []byte) error {
	type Alias TransactionResponse
	temp := &struct {
		Nonce            string `json:"nonce"`
		BlockNumber      string `json:"blockNumber"`
		TransactionIndex string `json:"transactionIndex"`
		Value            string `json:"value"`
		GasPrice         string `json:"gasPrice,omitempty"`
		Gas              string `json:"gas,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	nonce, err := utils.NewBigIntFromHex(temp.Nonce)

	if err != nil {
		return err
	}

	blockNum, err := utils.NewBigIntFromHex(temp.BlockNumber)

	if err != nil {
		return err
	}

	txIndex, err := utils.NewBigIntFromHex(temp.TransactionIndex)

	if err != nil {
		return err
	}

	gas, err := utils.NewBigIntFromHex(temp.Gas)

	if err != nil {
		return err
	}

	gasPrice, err := utils.NewBigIntFromHex(temp.GasPrice)

	if err != nil {
		return err
	}

	value, err := utils.NewBigIntFromHex(temp.Value)

	if err != nil {
		return err
	}

	t.Nonce = nonce
	t.BlockNumber = blockNum
	t.TransactionIndex = txIndex
	t.Gas = gas
	t.GasPrice = gasPrice
	t.Value = value

	return nil
}

func (r *TransactionLog) UnmarshalJSON(data []byte) error {
	type Alias TransactionLog

	log := &struct {
		TransactionIndex string `json:"transactionIndex"`
		BlockNumber      string `json:"blockNumber"`
		LogIndex         string `json:"logIndex"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &log); err != nil {
		return err
	}

	blockNumLog, err := utils.NewBigIntFromHex(log.BlockNumber)

	if err != nil {
		return err
	}

	txIndexLogs, err := utils.NewBigIntFromHex(log.TransactionIndex)

	if err != nil {
		return err
	}

	logIndex, err := utils.NewBigIntFromHex(log.LogIndex)

	if err != nil {
		return err
	}

	r.BlockNumber = blockNumLog
	r.TransactionIndex = txIndexLogs
	r.LogIndex = logIndex

	return nil

}

func (r *TransactionReceipt) UnmarshalJSON(data []byte) error {
	type Alias TransactionReceipt

	temp := &struct {
		TransactionIndex  string `json:"transactionIndex"`
		BlockNumber       string `json:"blockNumber"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		GasUsed           string `json:"gasUsed"`
		Status            string `json:"status"`
		Type              string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	blockNum, err := utils.NewBigIntFromHex(temp.BlockNumber)

	if err != nil {
		return err
	}

	txIndex, err := utils.NewBigIntFromHex(temp.TransactionIndex)

	if err != nil {
		return err
	}

	gasUsed, err := utils.NewBigIntFromHex(temp.GasUsed)

	if err != nil {
		return err
	}

	cumulativeGas, err := utils.NewBigIntFromHex(temp.CumulativeGasUsed)

	if err != nil {
		return err
	}

	status, err := utils.NewBigIntFromHex(temp.Status)
	if err != nil {
		return err
	}

	stype, err := utils.NewBigIntFromHex(temp.Type)
	if err != nil {
		return err
	}

	r.TransactionIndex = txIndex
	r.BlockNumber = blockNum
	r.CumulativeGasUsed = cumulativeGas
	r.GasUsed = gasUsed
	r.Status = false
	r.Type = stype
	if status.Cmp(big.NewInt(1)) == 0 {
		r.Status = true
	}

	return nil
}
