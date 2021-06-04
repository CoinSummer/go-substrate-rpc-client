package gsrpc

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
)

// //0
func TestSubstrateTransfer(t *testing.T) {
	// This sample shows how to create a transaction to make a transfer from one an account to another.

	// Instantiate the API
	// wss://acala-mandala.api.onfinality.io/public-ws
	//rpcUrl := `wss://rococo-1.acala.laminar.one`
	api, err := NewSubstrateAPI(`wss://rpc.polkadot.io`)
	//api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Create a call, transferring 12345 units to kkk
	// kkk1
	bob, err := types.NewMultiAddressFromHexAccountID("0x6a5809129023611c02ed12061f92c93e820a7baef1c39429db4ff797ca2ede5d")
	if err != nil {
		panic(err)
	}

	// 1 unit of transfer
	bal, ok := new(big.Int).SetString("10000000000", 10)
	if !ok {
		panic(fmt.Errorf("failed to convert balance"))
	}

	c, err := types.NewCall(meta, "Balances.transfer", bob, types.NewUCompact(bal))
	if err != nil {
		panic(err)
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	prikey := `xxx`
	keyringPair, err := signature.KeyringPairFromSecret(prikey, 42)
	fmt.Println("keyringPair.uri:", keyringPair.URI)

	if err != nil {
		panic(err)
	}

	fmt.Println("keyringPair:", keyringPair)

	key, err := types.CreateStorageKey(meta, "System", "Account", keyringPair.PublicKey, nil)
	//key, err := types.CreateStorageKey(meta, "System", "Account", alicePair.PublicKey, nil)
	if err != nil {
		panic(err)
	}

	var accountInfo types.AccountInfo
	ok, err = api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(err)
	}

	nonce := uint32(accountInfo.Nonce)
	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction using Alice's default account
	err = ext.Sign(keyringPair, o)
	if err != nil {
		panic(err)
	}

	// Send the extrinsic
	_, err = api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Balance transferred from Alice to Bob: %v\n", bal.String())
}
