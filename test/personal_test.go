package test

import (
	"testing"

	"github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/providers"
)

var personal *web3.Web3 = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))

func TestListAccountsAndNewAccount(t *testing.T) {
	account, err := personal.Personal.NewAccount("test")

	if err != nil {
		t.Error("error in new account request: ", err)
		t.Fail()
	}

	accounts, err := personal.Personal.ListAccounts()

	if err != nil {
		t.Error("error in list accounts request: ", err)
		t.Fail()
	}

	matched := false
	for _, v := range accounts {
		if v == account {
			matched = true
			break
		}
	}

	if !matched {
		t.Error("new address not found on list accounts request")
		t.Fail()
	}
}

func TestUnlockAccountAndSendTransaction(t *testing.T) {
	//TODO
}
