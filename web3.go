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
 * @file web3.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package web3

import (
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
	web3 "github.com/cellcycle/go-web3/web3"
)

// Coin - Ethereum value unity value
const (
	Coin float64 = 1000000000000000000
)

// Web3 - The Web3 Module
type Web3 struct {
	Provider providers.ProviderInterface
	Eth      *web3.Eth
	Net      *web3.Net
	Personal *web3.Personal
	Utils    *web3.Utils
}

// NewWeb3 - Web3 Module constructor to set the default provider, Eth, Net and Personal
func NewWeb3(provider providers.ProviderInterface) *Web3 {
	web3Client := new(Web3)
	web3Client.Provider = provider
	web3Client.Eth = web3.NewEth(provider)
	web3Client.Net = web3.NewNet(provider)
	web3Client.Personal = web3.NewPersonal(provider)
	web3Client.Utils = web3.NewUtils(provider)
	return web3Client
}

// ClientVersion - Returns the current client version.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_clientversion
// Parameters:
//    - none
// Returns:
// 	  - String - The current client version
func (web Web3) ClientVersion() (string, error) {

	pointer := &dto.RequestResult{}

	err := web.Provider.SendRequest(pointer, "web3_clientVersion", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}
