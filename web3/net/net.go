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
 * @file net.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package net

import (
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
)

type Net struct {
	provider providers.ProviderInterface
}

// NewNet - creates a new net instance
func NewNet(provider providers.ProviderInterface) *Net {
	net := new(Net)
	net.provider = provider
	return net
}

// Reference: https://eth.wiki/json-rpc/API#net_listening
func (net *Net) IsListening() (bool, error) {

	pointer := &dto.RequestResult{}

	err := net.provider.SendRequest(pointer, "net_listening", nil)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()

}

// Reference: https://eth.wiki/json-rpc/API#net_peercount
func (net *Net) GetPeerCount() (int64, error) {

	pointer := &dto.RequestResult{}

	err := net.provider.SendRequest(pointer, "net_peerCount", nil)

	if err != nil {
		return -1, err
	}

	return pointer.ToInt()

}

// Reference: https://eth.wiki/json-rpc/API#net_version
func (net *Net) GetVersion() (int64, error) {

	pointer := &dto.RequestResult{}

	err := net.provider.SendRequest(pointer, "net_version", nil)

	if err != nil {
		return -1, err
	}

	return pointer.ToInt()

}
