package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/conformal/sumrpcclient"
	"github.com/conformal/sumutil"
	"github.com/conformal/sumwire"
)

// This example demonstrates a connection to the bitcoin network
// by using websockets via sumd, use of notifications and an rpc
// call to getblockcount.
//
// Install and run sumd:
// 		$ go get github.com/conformal/sumd/...
//		$ sumd -u rpcuser -P rpcpass
//
// Install sumrpcclient:
// 		$ go get github.com/conformal/sumrpcclient
//
// Run this example:
// 		$ go run websocket-example.go
//
func main() {
	// Only override the handlers for notifications you care about.
	// Also note most of these handlers will only be called if you register
	// for notifications. See the documentation of the sumrpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := sumrpcclient.NotificationHandlers{
		OnBlockConnected: func(hash *sumwire.ShaHash, height int32) {
			log.Printf("Block connected: %v (%d)", hash, height)
		},
		OnBlockDisconnected: func(hash *sumwire.ShaHash, height int32) {
			log.Printf("Block disconnected: %v (%d)", hash, height)
		},
	}

	// Connect to local sumd RPC server using websockets.
	sumdHomeDir := sumutil.AppDataDir("sumd", false)
	certs, err := ioutil.ReadFile(filepath.Join(sumdHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &sumrpcclient.ConnConfig{
		Host:         "localhost:8334",
		Endpoint:     "ws",
		User:         "rpcuser",
		Pass:         "rpcpass",
		Certificates: certs,
	}
	client, err := sumrpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}

	// Register for block connect and disconnect notifications.
	if err := client.NotifyBlocks(); err != nil {
		log.Fatal(err)
	}
	log.Println("NotifyBlocks: Registration Complete")

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)

	// For this example gracefully shutdown the client after 10 seconds.
	// Ordinarily when to shutdown the client is highly application
	// specific.
	log.Println("Client shutdown in 10 seconds...")
	time.AfterFunc(time.Second*10, func() {
		log.Println("Client shutting down...")
		client.Shutdown()
		log.Println("Client shutdown complete.")
	})

	// Wait until the client either shuts down gracefully (or the user
	// terminates the process with Ctrl+C).
	client.WaitForShutdown()
}
