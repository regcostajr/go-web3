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
 * @file websocket-provider.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package providers

import (
	"encoding/json"
	"errors"
	customerror "github.com/cellcycle/go-web3/constants"
	"github.com/cellcycle/go-web3/dto"
	"math/rand"
	"strings"

	"github.com/gorilla/websocket"
)

type WebSocketProvider struct {
	address string
	secure  bool
	ws      *websocket.Conn
}

func NewWebSocketProvider(address string, secure bool) *WebSocketProvider {
	provider := new(WebSocketProvider)
	provider.address = address
	provider.secure = secure
	return provider
}

func (provider *WebSocketProvider) Connect() error {
	if provider.ws == nil {

		prefix := ""
		if !strings.HasPrefix(provider.address, "ws") {
			prefix = "ws://"
			if provider.secure {
				prefix = "wss://"
			}
		}

		url := prefix + provider.address

		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			return err
		}
		provider.ws = ws
	}

	return nil
}

func (provider *WebSocketProvider) SendRequest(method string, params interface{}) (*dto.RequestResult, error) {
	if provider.ws == nil {
		return nil, errors.New("connection is closed")
	}

	bodyString := JSONRPCObject{Version: "2.0", Method: method, Params: params, ID: rand.Intn(100)}
	message := []byte(bodyString.AsJsonString())
	err := provider.ws.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return nil, err
	}

	_, response, err := provider.ws.ReadMessage()

	var requestResult *dto.RequestResult
	json.Unmarshal(response, &requestResult)

	return requestResult, err
}

func (provider *WebSocketProvider) Subscribe(ch chan<- *dto.Subscription, method string, params interface{}) error {
	if provider.ws == nil {
		return errors.New("connection is closed")
	}

	bodyString := JSONRPCObject{Version: "2.0", Method: method, Params: params, ID: rand.Intn(100)}
	message := []byte(bodyString.AsJsonString())
	err := provider.ws.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	go func() {
		defer close(ch)

		for {
			_, message, err := provider.ws.ReadMessage()
			if err != nil {
				return
			}

			var subscription *dto.Subscription
			json.Unmarshal(message, &subscription)

			ch <- subscription
		}
	}()

	return err
}

func (provider *WebSocketProvider) Close() error {
	if provider.ws != nil {
		err := provider.ws.Close()
		provider.ws = nil
		return err
	}

	return customerror.WEBSOCKETNOTDENIFIED
}
