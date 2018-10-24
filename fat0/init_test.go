package fat0_test

import (
	"encoding/hex"
	"encoding/json"
	"math/rand"

	"github.com/Factom-Asset-Tokens/fatd/factom"
	"github.com/Factom-Asset-Tokens/fatd/fat0"
	"github.com/FactomProject/ed25519"
)

const (
	blockheight uint64 = 160000
)

var (
	validIssuerChainID *factom.Bytes32
	issuerKey          factom.Address

	validIssuanceEntryContentMap = map[string]interface{}{
		"type":   "FAT-0",
		"supply": int64(100000),
		"symbol": "TEST",
		"name":   "Test Token",
	}
	validIssuanceEntry factom.Entry
	validIssuance      *fat0.Issuance

	validTransactionEntryContentMap map[string]interface{}
	validTransactionEntry           factom.Entry
	validTransaction                *fat0.Transaction
	inputs                          []factom.Address
	inputAmounts                    = []uint64{100, 10}
	outputs                         []factom.Address
	outputAmounts                   = []uint64{90, 20}
)

type addressAmount struct {
	Address factom.Address `json:"address"`
	Amount  uint64         `json:"amount"`
}

func init() {
	id, _ := hex.DecodeString(
		"88888807e4f3bbb9a2b229645ab6d2f184224190f83e78761674c2362aca4425")
	validIssuerChainID = factom.NewBytes32(id)
	validIssuanceEntry.EBlock = factom.EBlock{
		ChainID: fat0.ChainID("test", validIssuerChainID)}
	validIssuanceEntry.Content = marshal(validIssuanceEntryContentMap)

	rand := rand.New(rand.NewSource(100))
	issuerKey.PublicKey, issuerKey.PrivateKey, _ = ed25519.GenerateKey(rand)
	validIssuance = fat0.NewIssuance(&validIssuanceEntry)
	validIssuance.Sign(issuerKey)

	inputs = make([]factom.Address, 2)
	outputs = make([]factom.Address, 2)
	for _, addresses := range [][]factom.Address{inputs, outputs} {
		for i := range addresses {
			addresses[i].PublicKey, addresses[i].PrivateKey, _ =
				ed25519.GenerateKey(rand)
		}
	}

	validTransactionEntryContentMap = map[string]interface{}{
		"inputs": []addressAmount{{
			Address: inputs[0],
			Amount:  inputAmounts[0],
		}, {
			Address: inputs[1],
			Amount:  inputAmounts[1],
		}},
		"outputs": []addressAmount{{
			Address: outputs[0],
			Amount:  outputAmounts[0],
		}, {
			Address: outputs[1],
			Amount:  outputAmounts[1],
		}},
		"blockheight": blockheight,
		"salt":        "xyz",
	}

	validTransactionEntry.Content = marshal(validTransactionEntryContentMap)
	validTransactionEntry.ChainID = validIssuanceEntry.ChainID
	validTransactionEntry.Height = blockheight

	validTransaction = fat0.NewTransaction(&validTransactionEntry)
	validTransaction.Sign(inputs...)
}

func marshal(v map[string]interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

func mapCopy(dst, src map[string]interface{}) {
	for k, v := range src {
		dst[k] = v
	}
}
