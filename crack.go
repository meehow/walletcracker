package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var numWorkers = runtime.NumCPU()

func crack(addresses []common.Address) {
	keyChan := make(chan *ecdsa.PrivateKey)
	for i := 0; i < numWorkers; i++ {
		go generateWallet(keyChan)
	}
	counter := 0
	start := time.Now()
	for privateKey := range keyChan {
		if counter%100000 == 0 {
			rate := float64(counter) / time.Since(start).Seconds()
			fmt.Printf("Keys per second: %.0f\r", rate)
		}
		counter++
		address := crypto.PubkeyToAddress(privateKey.PublicKey)
		if contains(addresses, address) {
			privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
			fmt.Printf("Address: %s PrivateKey: %s\n", address.Hex(), privateKeyHex)
		}
	}
}

func contains(addresses []common.Address, address common.Address) bool {
	// addresses must be sorted to make binary search working
	i, j := 0, len(addresses)
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		cmp := bytes.Compare(addresses[h][:], address[:])
		if cmp < 0 {
			i = h + 1
		} else if cmp > 0 {
			j = h
		} else {
			return true
		}
	}
	return false
}

func generateWallet(keyChan chan *ecdsa.PrivateKey) {
	for {
		privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
		if err != nil {
			log.Fatal(err)
		}
		keyChan <- privateKey
	}
}
