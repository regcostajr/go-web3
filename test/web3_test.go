package test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/cellcycle/go-web3"
	"github.com/cellcycle/go-web3/providers"
)

var web3Client *web3.Web3 = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))

func TestSha3(t *testing.T) {
	cipherText := web3Client.Sha3("test")
	// curl 127.0.0.1:8545 \
	// -X POST \
	// -H "Content-Type: application/json" \
	// -d '{"jsonrpc":"2.0","method":"web3_sha3","params":["0x74657374"],"id":1}'
	test_result := "0x9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658"

	if cipherText != test_result {
		t.Error("wrong sha3 conversion")
		t.Fail()
	}

	cipherText = web3Client.Sha3("Transfer(address,address,uint256)")
	// https://etherscan.io/tx/0xe317330e165b27df0728e53bb605b288f550b0a8c1e46146699c914d5a447aa6#eventlog
	test_result = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

	if cipherText != test_result {
		t.Error("wrong sha3 conversion")
		t.Fail()
	}
}

func TestClientVersion(t *testing.T) {
	version, err := web3Client.ClientVersion()

	if err != nil || version == "" {
		t.Error("wrong response for client version", err)
		t.Fail()
	}

	if strings.HasPrefix(version, "Geth") {
		match, err := regexp.Match(`^Geth\/v\d+.\d+.\d+-[a-z]+-[a-z0-9]+`, []byte(version))
		if match && err != nil {
			t.Error("wrong response for client version", err)
			t.Fail()
		}
	}
}
