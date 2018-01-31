package main

import (
	"log"
	"strconv"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	bc := NewBlockChain()
	bc.AddBlock("Send 1 BTC to Elvin")
	bc.AddBlock("Send 2 more BTC to Elvin")

	for _, block := range bc.blocks {
		log.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		log.Printf("Data. %s\n", block.Data)
		log.Printf("Hash. %x\n", block.Hash)

		pow := NewProofOfWork(block)
		log.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	}
}
