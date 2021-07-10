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
	"fmt"

	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
	"github.com/cellcycle/go-web3/web3/eth"
	"github.com/cellcycle/go-web3/web3/net"
	"github.com/cellcycle/go-web3/web3/personal"
	"golang.org/x/crypto/sha3"
)

type Web3 struct {
	Provider providers.ProviderInterface
	Eth      *eth.Eth
	Net      *net.Net
	Personal *personal.Personal
}

// NewWeb3 - creates a new web3 instance
func NewWeb3(provider providers.ProviderInterface) *Web3 {
	web3Client := new(Web3)
	web3Client.Provider = provider
	web3Client.Eth = eth.NewEth(provider)
	web3Client.Net = net.NewNet(provider)
	web3Client.Personal = personal.NewPersonal(provider)
	return web3Client
}

// Reference: https://eth.wiki/json-rpc/API#web3_clientversion
func (web3 *Web3) ClientVersion() (string, error) {
	pointer := &dto.RequestResult{}

	err := web3.Provider.SendRequest(pointer, "web3_clientVersion", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

// Sha3 - hashes the given string using keccak-256 and returns
// an hexadecimal hash containing the prefix 0x
func (web3 *Web3) Sha3(data string) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(data))
	return fmt.Sprintf("0x%x", hash.Sum(nil))
}
