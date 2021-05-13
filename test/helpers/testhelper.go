package helpers

import (
	"github.com/cellcycle/go-web3/dto"
	"math/big"
)

func GetTestTx(to string, from string, amount *big.Int) *dto.TransactionParameters {

	Tx := new(dto.TransactionParameters)

	transaction := new(dto.TransactionParameters)
	transaction.From = from
	transaction.To = to
	transaction.Value = amount
	transaction.Gas = big.NewInt(40000)

	return Tx
}
