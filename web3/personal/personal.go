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
 * @file personal.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package personal

import (
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
)

type Personal struct {
	provider providers.ProviderInterface
}

// NewPersonal - creates a new personal instance
func NewPersonal(provider providers.ProviderInterface) *Personal {
	personal := new(Personal)
	personal.provider = provider
	return personal
}

// Reference: https://geth.ethereum.org/docs/rpc/ns-personal#personal_listaccounts
func (personal *Personal) ListAccounts() ([]string, error) {
	pointer := &dto.RequestResult{}

	err := personal.provider.SendRequest(pointer, "personal_listAccounts", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToStringArray()
}

// Reference: https://geth.ethereum.org/docs/rpc/ns-personal#personal_newaccount
func (personal *Personal) NewAccount(password string) (string, error) {

	params := make([]string, 1)
	params[0] = password

	pointer := &dto.RequestResult{}

	err := personal.provider.SendRequest(&pointer, "personal_newAccount", params)

	if err != nil {
		return "", err
	}

	response, err := pointer.ToString()

	return response, err

}

// Reference: https://github.com/paritytech/parity/wiki/JSONRPC-personal-module#personal_sendtransaction
func (personal *Personal) SendTransaction(transaction *dto.TransactionParameters, password string) (string, error) {

	params := make([]interface{}, 2)

	transactionParameters := transaction.Transform()

	params[0] = transactionParameters
	params[1] = password

	pointer := &dto.RequestResult{}

	err := personal.provider.SendRequest(pointer, "personal_sendTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

// Reference: https://geth.ethereum.org/docs/rpc/ns-personal#personal_unlockaccount
func (personal *Personal) UnlockAccount(address string, password string, duration uint64) (bool, error) {

	params := make([]interface{}, 3)
	params[0] = address
	params[1] = password
	params[2] = duration

	pointer := &dto.RequestResult{}

	err := personal.provider.SendRequest(pointer, "personal_unlockAccount", params)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}
