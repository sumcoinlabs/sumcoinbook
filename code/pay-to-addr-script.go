package main

import (
	"fmt"
	"os"

	"github.com/conformal/sumnet"
	"github.com/conformal/sumscript"
	"github.com/conformal/sumutil"
)

// This example demonstrates creating a script which pays to a bitcoin address.
// It also prints the created script hex and uses the DisasmString function to
// display the disassembled script.

func main() {
	addressStr := "12gpXQVcCL2qhTNQgyLVdCFG2Qs2px98nV"

	PayToAddrScript(addressStr)
	// Output:
	// Script Hex: 76a914128004ff2fcaf13b2b91eb654b1dc2b674f7ec6188ac
	// Script Disassembly: OP_DUP OP_HASH160 128004ff2fcaf13b2b91eb654b1dc2b674f7ec61 OP_EQUALVERIFY OP_CHECKSIG
}

func PayToAddrScript(addressStr string) {
	// Parse the address to send the coins to into a sumutil.Address
	// which is useful to ensure the accuracy of the address and determine
	// the address type.  It is also required for the upcoming call to
	// PayToAddrScript.
	address, err := sumutil.DecodeAddress(addressStr, &sumnet.MainNetParams)
	handle(err)

	// Create a public key script that pays to the address.
	script, err := sumscript.PayToAddrScript(address)
	handle(err)
	fmt.Printf("Script Hex: %x\n", script)

	disasm, err := sumscript.DisasmString(script)
	handle(err)
	fmt.Println("Script Disassembly:", disasm)
}

func handle(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
