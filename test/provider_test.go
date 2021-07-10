package test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
)

func Test_HttpProvider(t *testing.T) {

	var provider = providers.NewHTTPProvider("127.0.0.1:8545", 10, false)

	var pointer interface{}
	err := provider.SendRequest(&pointer, "net_listening", nil)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	response := pointer.(map[string]interface{})

	if !response[`result`].(bool) {
		t.Error("wrong response for listening node")
		t.Fail()
	}

}

func Test_WebSocketProviderSimpleRequest(t *testing.T) {

	var provider = providers.NewWebSocketProvider("127.0.0.1:8546", false)
	err := provider.Connect()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	message, err := provider.SendRequest("eth_blockNumber", nil)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	block, err := message.ToInt()

	if err != nil || block <= 0 {
		t.Error(err)
		t.Fail()
	}

	provider.Close()

	message, err = provider.SendRequest("eth_blockNumber", nil)

	if err == nil {
		t.Error("using provider after connection closed, it should thrown an error")
		t.Fail()
	}

	// test reconnection
	err = provider.Connect()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	message, err = provider.SendRequest("eth_blockNumber", nil)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	block, err = message.ToInt()

	if err != nil || block <= 0 {
		t.Error(err)
		t.Fail()
	}

}

func Test_WebSocketProviderSubscription(t *testing.T) {
	infuraConfig, err := ioutil.ReadFile("../infura.conf")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	infuraToken := strings.TrimSuffix(string(infuraConfig), "\n")

	var provider = providers.NewWebSocketProvider("mainnet.infura.io/ws/v3/"+infuraToken, true)
	err = provider.Connect()
	defer provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	params := make([]string, 1)
	params[0] = "newPendingTransactions"

	channel := make(chan *dto.Subscription)
	err = provider.Subscribe(channel, "eth_subscribe", params)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	for i := 0; i < 11; i++ {
		response := <-channel

		// first response is the subscription confirmation
		if response.ID > 0 {
			continue
		}
		if len(response.Params.Result) < 66 {
			t.Error("transaction size does not matches")
			t.Fail()
		}
	}
}
