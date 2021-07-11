package dto

import (
	"encoding/hex"
	"fmt"
	"github.com/cellcycle/go-web3/utils"
	"math/big"
)

type DataChunks []string

func (chunks DataChunks) DecodeInt(index int64) *big.Int {
	big := new(big.Int)
	big.SetString(chunks[index], 16)
	return big
}

func (chunks DataChunks) DecodeString(index int64) (string, error) {
	data, err := hex.DecodeString(chunks[index])
	if err != nil {
		return "", err
	}

	return utils.CleanString(string(data)), nil
}

func (chunks DataChunks) DecodeBytes(index int64) ([]byte, error) {
	strBytes, err := chunks.DecodeString(index)

	if err != nil {
		return nil, err
	}

	return []byte(strBytes), nil
}

func (chunks DataChunks) DecodeAddress(index int64) string {
	address := chunks[index]
	return utils.CleanString(fmt.Sprintf("0x%s", address[len(address)-40:]))
}

func (chunks DataChunks) ToString(index int64) string {
	return string(chunks[index])
}
