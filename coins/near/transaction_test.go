package near

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/okx/go-wallet-sdk/coins/near/serialize"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCreateTransaction(t *testing.T) {

	privateKey := "b9ec4d26ab5bec8df4314a9e3b8fc3f9c96d410b4cd13caa675018dcfc7916cceefbba85caaa14cb87b83314d5b86895f2d4b7633e29012e65bfb037c885c804"
	val := 0.222
	to := "ggasii.testnet"
	blockHash := "D7CPxgTXyRKYTSYuwAiRwDJH5RdHz7DwPt4EViptAW4L"
	nonce := int64(1)

	addr, err := PrivateKeyToAddr(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	publicKeyHex, err := PrivateKeyToPublicKeyHex(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := CreateTransaction(addr, to, publicKeyHex, blockHash, nonce)
	if err != nil {
		t.Fatal(err)
	}

	amount := decimal.NewFromFloat(val).Shift(24)
	ta, err := serialize.CreateTransfer(amount.String())
	if err != nil {
		t.Fatal(err)
	}
	tx.SetAction(ta)
	txData, err := tx.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	txBase58 := base58.Encode(txData)
	sig, err := SignTransaction(txBase58, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	stx, err := CreateSignedTransaction(tx, sig)
	if err != nil {
		t.Fatal(err)
	}
	stxData, err := stx.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	b64Data := base64.StdEncoding.EncodeToString(stxData)
	t.Logf("transaction data : %s", b64Data)
	assert.Equal(t, "QAAAAGVlZmJiYTg1Y2FhYTE0Y2I4N2I4MzMxNGQ1Yjg2ODk1ZjJkNGI3NjMzZTI5MDEyZTY1YmZiMDM3Yzg4NWM4MDQA7vu6hcqqFMuHuDMU1bholfLUt2M+KQEuZb+wN8iFyAQBAAAAAAAAAA4AAABnZ2FzaWkudGVzdG5ldLPinQIWXUUnnN9Qmtou83BpsylI4Fb+ZStWsef3s/kNAQAAAAMAAMAOl7HkpAIvAAAAAAAAACE/E/jQF9vlZSvRNf3Dnrr9Tm+gPB4s4wvE46LM18lgPtighyOfczJQMwhTJjFBe5xzBWbq3CJVhUYK21a9nQ0=", b64Data)
}

func TestContactTransaction(t *testing.T) {

	privateKey := "b9ec4d26ab5bec8df4314a9e3b8fc3f9c96d410b4cd13caa675018dcfc7916cceefbba85caaa14cb87b83314d5b86895f2d4b7633e29012e65bfb037c885c804"
	val := 0.222
	to := "ft.examples.testnet"
	blockHash := "D7CPxgTXyRKYTSYuwAiRwDJH5RdHz7DwPt4EViptAW4L"
	nonce := int64(1)

	argsStr := `{"account_id": "serhii.testnet"}`
	gas := big.NewInt(1)

	addr := "ggasii.testnet"

	publicKeyHex, err := PrivateKeyToPublicKeyHex(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := CreateTransaction(addr, to, publicKeyHex, blockHash, nonce)
	if err != nil {
		t.Fatal(err)
	}

	amount := decimal.NewFromFloat(val).Shift(24).BigInt()
	println(amount.String())
	ta, err := serialize.CreateFunctionCall("storage_balance_of", []byte(argsStr), gas, amount)
	if err != nil {
		t.Fatal(err)
	}

	tx.SetAction(ta)
	txData, err := tx.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	txHash := base58.Encode(txData)
	sig, err := SignTransaction(txHash, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	stx, err := CreateSignedTransaction(tx, sig)
	if err != nil {
		t.Fatal(err)
	}
	stxData, err := stx.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	println(hex.EncodeToString(stxData))

	b64Data := base64.StdEncoding.EncodeToString(stxData)
	t.Logf("transaction data : %s", b64Data)
	assert.Equal(t, "DgAAAGdnYXNpaS50ZXN0bmV0AO77uoXKqhTLh7gzFNW4aJXy1LdjPikBLmW/sDfIhcgEAQAAAAAAAAATAAAAZnQuZXhhbXBsZXMudGVzdG5ldLPinQIWXUUnnN9Qmtou83BpsylI4Fb+ZStWsef3s/kNAQAAAAISAAAAc3RvcmFnZV9iYWxhbmNlX29mIAAAAHsiYWNjb3VudF9pZCI6ICJzZXJoaWkudGVzdG5ldCJ9AQAAAAAAAAAAAMAOl7HkpAIvAAAAAAAAACvEiv+vj1JDfHnrGZZ9vQlVvKCb2Bqsqe2KBB3ZhyM1YcWRR6WvjWVWpBmiXHt48xUf8ePtVcKdc0BNau8bJQM=", b64Data)
}
