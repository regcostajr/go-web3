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
 * @file db.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package db

import (
	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
)

// DB - The DB Module
type DB struct {
	provider providers.ProviderInterface
}

// NewDB - DB Module constructor to set the default provider
func NewDB(provider providers.ProviderInterface) *DB {
	db := new(DB)
	db.provider = provider
	return db
}

// PutString - Stores a string in the local database.
// Note this function is deprecated and will be removed in the future.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#db_putstring
// Parameters:
//    - String - Database name.
//	  - String - Key name.
//    - String - String to store.
// Returns:
//	  - Boolean - returns true if the value was stored, otherwise false.
func (db *DB) PutString(databaseName string, keyName string, stringToStore string) (bool, error) {

	params := make([]string, 3)

	params[0] = databaseName
	params[1] = keyName
	params[2] = stringToStore

	pointer := &dto.RequestResult{}

	err := db.provider.SendRequest(pointer, "db_putString", params)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()

}
