package test

import (
	"testing"

	"github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/providers"
)

var net *web3.Web3 = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))

func TestIsListening(t *testing.T) {
	isListening, err := net.Net.IsListening()

	if !isListening || err != nil {
		t.Error("wrong response for listening node", err)
		t.Fail()
	}
}

func TestGetPeerCount(t *testing.T) {
	count, err := net.Net.GetPeerCount()

	if count == -1 || err != nil {
		t.Error("wrong response for peer count", err)
		t.Fail()
	}
}

func TestGetVersion(t *testing.T) {
	version, err := net.Net.GetVersion()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if version != 4919 {
		t.Error("wrong version for dev mode geth, received: ", version)
		t.Fail()
	}
}
