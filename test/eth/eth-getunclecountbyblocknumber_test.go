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
 * @file eth-getunclecountbyblocknumber_test.go
 * @authors:
 * 		Sigma Prime <sigmaprime.io>
 * @date 2018
 */

package test

import (
	"github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/providers"
	"math/big"
	"testing"
)

func TestGetUncleCountByBlockNumber(t *testing.T) {

	var connection = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))

	blockNumber, err := connection.Eth.GetBlockNumber()

	uncleByNumber, err := connection.Eth.GetUncleCountByBlockNumber(blockNumber)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(uncleByNumber.Int64())

	if uncleByNumber.Int64() != 0 {
		t.Errorf("Returned uncle for block with no uncle.")
		t.FailNow()
	}

	// should return err with negative number?
	uncleByNumber, err = connection.Eth.GetUncleCountByBlockNumber(big.NewInt(-1))

	if err == nil {
		t.Error(err)
		t.FailNow()
	}
}
