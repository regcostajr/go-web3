package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cellcycle/go-web3/dto"
	"github.com/cellcycle/go-web3/providers"
)

func main() {
	infuraConfig, err := ioutil.ReadFile("infura.conf")
	if err != nil {
		fmt.Println(err)
	}

	infuraToken := strings.TrimSuffix(string(infuraConfig), "\n")

	var provider = providers.NewWebSocketProvider("mainnet.infura.io/ws/v3/"+infuraToken, true)
	err = provider.Connect()
	defer provider.Close()

	if err != nil {
		fmt.Println(err)
	}

	params := make([]string, 1)
	params[0] = "newPendingTransactions"

	channel := make(chan *dto.Subscription)
	err = provider.Subscribe(channel, "eth_subscribe", params)

	if err != nil {
		fmt.Println(err)
	}

	for {
		response := <-channel
		fmt.Println(response.Params.Result)
	}
}
