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
 * @file eth-getblocktransactioncountbynumber.go
 * @authors:
 *   Junjie CHen <chuckjunjchen@gmail.com>
 * @date 2018
 */

package test

import (
	"math/big"
	"testing"
	"time"

	web3 "github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/eth/block"
	"github.com/cellcycle/go-web3/providers"
)

func TestGetBlockTransactionCountByNumber(t *testing.T) {

	var connection = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))

	// submit a transaction, wait for the block and there should be 1 tx.
	coinbase, err := connection.Eth.GetCoinbase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = coinbase
	transaction.To = coinbase
	transaction.Value = big.NewInt(200000)
	transaction.Gas = big.NewInt(40000)

	txID, err := connection.Eth.SendTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	time.Sleep(time.Second)

	tx, err := connection.Eth.GetTransactionByHash(txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	blockNumber := block.NUMBER(tx.BlockNumber)

	txCount, err := connection.Eth.GetBlockTransactionCountByNumber(blockNumber)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if txCount.Int64() != 1 {
		t.Error("invalid block transaction count")
		t.FailNow()
	}

	txCount, err = connection.Eth.GetBlockTransactionCountByNumber(block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if txCount.Int64() != 1 {
		t.Error("invalid block transaction count")
		t.FailNow()
	}
}
