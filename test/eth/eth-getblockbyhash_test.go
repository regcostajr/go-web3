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
 * @file eth-getBlockByNumber_test.go
 * @authors:
 *    Sigma Prime <sigmaprime.io>
 * @date 2018
 */

package test

import (
	"testing"

	"github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/providers"
)

func TestEthGetBlockByHash(t *testing.T) {

	var connection = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))

	blockNumber, err := connection.Eth.GetBlockNumber()

	blockByNumber, err := connection.Eth.GetBlockByNumber(blockNumber, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	blockByHash, err := connection.Eth.GetBlockByHash(blockByNumber.Hash, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Ensure it's the same block
	if (blockByNumber.Number.Cmp(blockByHash.Number)) != 0 ||
		(blockByNumber.Miner != blockByHash.Miner) ||
		(blockByNumber.Hash != blockByHash.Hash) {
		t.Errorf("Not same block returned")
		t.FailNow()
		t.FailNow()
	}

	t.Log(blockByHash.Hash, blockByNumber.Hash)

	_, err = connection.Eth.GetBlockByHash("0x1234", false)

	if err == nil {
		t.Errorf("Invalid hash not rejected")
		t.FailNow()
	}

	_, err = connection.Eth.GetBlockByHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa0", false)

	if err == nil {
		t.Errorf("Invalid hash not rejected")
		t.FailNow()
	}

	_, err = connection.Eth.GetBlockByHash("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa0", false)

	if err == nil {
		t.Errorf("Invalid hash not rejected")
		t.FailNow()
	}

	blockByHash, err = connection.Eth.GetBlockByHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", false)

	if err == nil {
		t.Errorf("Found a block with incorrect hash?")
		t.FailNow()
	}
}
