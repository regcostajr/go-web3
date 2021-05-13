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
 * @file shh.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package shh

import (
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
	"github.com/cellcycle/go-web3/utils"
	"math/big"
)

// SHH - The Net Module
type SHH struct {
	provider providers.ProviderInterface
}

// NewSHH - Net Module constructor to set the default provider
func NewSHH(provider providers.ProviderInterface) *SHH {
	shh := new(SHH)
	shh.provider = provider
	return shh
}

// GetVersion - Returns the current whisper protocol version.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#shh_version
// Parameters:
//    - none
// Returns:
// 	  - String - The current whisper protocol version
func (shh *SHH) GetVersion() (string, error) {

	pointer := &dto.RequestResult{}

	err := shh.provider.SendRequest(pointer, "shh_version", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

// Post - Sends a whisper message.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#shh_post
// Parameters:
//     1 .Object - The whisper post object:
//	  	- from: DATA, 60 Bytes - (optional) The identity of the sender.
//    	- to: DATA, 60 Bytes - (optional) The identity of the receiver. When present whisper will encrypt the message so that only the receiver can decrypt it.
//   	- topics: Array of DATA - Array of DATA topics, for the receiver to identify messages.
//    	- payload: DATA - The payload of the message.
//    	- priority: QUANTITY - The integer of the priority in a rang from ... (?).
//    	- ttl: QUANTITY - integer of the time to live in seconds.
// Returns:
// 	  - Boolean - returns true if the message was send, otherwise false.
func (shh *SHH) Post(from string, to string, topics []string, payload string, priority *big.Int, ttl *big.Int) (bool, error) {

	params := make([]dto.SHHPostParameters, 1)
	params[0].From = from
	params[0].To = to
	params[0].Topics = topics
	params[0].Payload = payload
	params[0].Priority = utils.IntToHex(priority)
	params[0].TTL = utils.IntToHex(ttl)

	pointer := &dto.RequestResult{}

	err := shh.provider.SendRequest(pointer, "shh_post", params)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()

}
