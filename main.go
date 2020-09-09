package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

const fName = "addresses.json"

func main() {
	addresses := []common.Address{}
	file, err := os.Open(fName)
	if err != nil {
		log.Printf("%v: crawling", err)
		file, err = os.Create(fName)
		if err != nil {
			log.Fatalln(err)
		}
		addresses = crawl()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "\t")
		err = encoder.Encode(addresses)
	} else {
		err = json.NewDecoder(file).Decode(&addresses)
	}
	file.Close()
	if err != nil {
		log.Fatalln(err)
	}
	crack(addresses)
}
