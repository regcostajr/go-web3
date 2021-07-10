package utils

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func NewBigIntFromHex(hexString string) (*big.Int, error) {
	if hexString == "0x" || hexString == "" {
		return big.NewInt(0), nil
	}

	bigInt, sucess := big.NewInt(0).SetString(strings.TrimPrefix(hexString, "0x"), 16)

	if !sucess {
		return nil, errors.New(fmt.Sprintf("can't convert hex: %s to big.Int", hexString))
	}

	return bigInt, nil
}

func CleanString(str string) string {
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
