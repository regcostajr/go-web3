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
 * @file complex-string.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package types

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type String string

func NewString(str string) *String {
	strInstance := String(str)
	return &strInstance
}

func (s String) ToHex() *String {
	if strings.HasPrefix(string(s), "0x") {
		return &s
	}
	return NewString(fmt.Sprintf("0x%x", s.ToString()))
}

func (s String) ToString() string {
	stringValue := string(s)
	if !strings.HasPrefix(stringValue, "0x") {
		return stringValue
	}

	sResult, _ := hex.DecodeString(strings.TrimPrefix(stringValue, "0x"))

	return s.clean(string(sResult))

}

func (s String) ToBytes() []byte {
	return []byte(s)
}

func (s String) Add0x() *String {
	if !strings.HasPrefix(string(s), "0x") {
		s = String("0x" + string(s))
	}
	return &s
}

func (s String) Concat(str string) *String {
	s = String(string(s) + str)
	return &s
}

func (s String) clean(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return strings.TrimSpace(string(b[:bl]))
}

func (s String) IsEmpty() bool {
	return string(s) == ""
}
