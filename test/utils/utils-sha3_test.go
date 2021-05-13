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
 * @file web3-sha3_test.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package test

import (
	"errors"
	"strings"
	"testing"

	web3 "github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/providers"
)

func TestUtilsSha3(t *testing.T) {

	var connection = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))

	sha3String, err := connection.Utils.Sha3("test")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(sha3String)

	result := strings.Compare("0x9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658", sha3String)

	if result != 0 {
		t.Error(errors.New("Sha3 string not equal calculed hash"))
		t.Fail()
	}

}
