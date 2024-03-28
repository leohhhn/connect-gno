package main

import (
	"fmt"
	"github.com/gnolang/gno/gno.land/pkg/gnoclient"
	rpcclient "github.com/gnolang/gno/tm2/pkg/bft/rpc/client"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys"
)

func main() {
	// Initialize keybase from a directory
	keybase, _ := keys.NewKeyBaseFromDir("path/to/keybase/dir")

	// Create signer
	signer := gnoclient.SignerFromKeybase{
		Keybase:  keybase,
		Account:  "<keypair_name>",     // Name of your keypair in keybase
		Password: "<keypair_password>", // Password to decrypt your keypair
		ChainID:  "<gno_chainID>",      // id of Gno.land chain
	}

	// Initialize the RPC client
	rpc := rpcclient.NewHTTP("<gno.land_remote_endpoint>", "")

	// Initialize the gnoclient
	client := gnoclient.Client{
		Signer:    signer,
		RPCClient: rpc,
	}

	// Convert Gno address string to `crypto.Address`
	addr, err := crypto.AddressFromBech32("g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5") // your Gno address
	if err != nil {
		panic(err)
	}

	// Get account info
	accountRes, _, err := client.QueryAccount(addr)
	if err != nil {
		panic(err)
	}

	// Construct base transaction config set
	txCfg := gnoclient.BaseTxCfg{
		GasFee:         "1000000ugnot",                 // gas price
		GasWanted:      1000000,                        // gas limit
		AccountNumber:  accountRes.GetAccountNumber(),  // account ID
		SequenceNumber: accountRes.GetSequence(),       // account nonce
		Memo:           "This is a cool how-to guide!", // transaction memo
	}

	// Construct message to pack into transaction
	msg := gnoclient.MsgCall{
		PkgPath:  "gno.land/r/demo/wugnot", // wrapped ugnot realm path
		FuncName: "Deposit",                // function to call
		Args:     nil,                      // arguments in string format
		Send:     "1000000ugnot",           // coins to send along with transaction
	}

	// Send transaction to chain
	res, err := client.Call(txCfg, msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	// Evaluate expression
	qevalRes, _, _ := client.QEval("gno.land/r/demo/wugnot", "BalanceOf(\"g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5\")")
	fmt.Println(qevalRes)

	// Get render output of a realm
	renderRes, _, _ := client.Render("gno.land/r/demo/echo", "Echo this text!")
	fmt.Println(renderRes)
}
