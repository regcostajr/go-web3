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
 * @file ipc-provider.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package providers

import (
	"encoding/json"
	"math/rand"
	"net"
	"path/filepath"

	"log"

	"github.com/cellcycle/go-web3/providers/util"
)

type IPCProvider struct {
	endpoint string
}

func NewIPCProvider(endpoint string) *IPCProvider {
	provider := new(IPCProvider)
	provider.endpoint, _ = filepath.Abs(endpoint)
	return provider
}

func (provider IPCProvider) SendRequest(v interface{}, method string, params interface{}) error {

	bodyString := util.JSONRPCObject{Version: "2.0", Method: method, Params: params, ID: rand.Intn(100)}

	client, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: provider.endpoint, Net: "unix"})

	if err != nil {
		log.Println(err)
		return err
	}

	defer client.Close()

	encoder := json.NewEncoder(client)
	decoder := json.NewDecoder(client)

	if err := encoder.Encode(bodyString); err != nil {
		log.Println(err)
		return err
	}

	if err := decoder.Decode(v); err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func (provider IPCProvider) Close() error { return nil }
