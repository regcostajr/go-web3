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
 * @file websocket-provider_test.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2018
 */
package providers

import (
	"testing"

	web3 "github.com/cellcycle/go-web3"
)

func Test_WebSocketProvider(t *testing.T) {

	var ethClient = web3.NewWeb3(providers.NewWebSocketProvider("ws://127.0.0.1:8546"))

	for index := 0; index < 100; index++ {

		var _, error = ethClient.ClientVersion()

		if error != nil {
			t.Error(error)
			t.Fail()
		}

	}

	ethClient.Provider.Close()

}
