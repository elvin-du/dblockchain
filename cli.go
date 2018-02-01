package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	bc *BlockChain
}

func (cli *CLI) Run() {
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block Data")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if nil != err {
			log.Printf("parse add block command failed,err:%s", err.Error())
			return
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if nil != err {
			log.Printf("parse print chain command failed,err:%s", err.Error())
			return
		}
	default:
		cli.printUseage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if "" == *addBlockData {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()
	for {
		block, err := bci.Next()
		if nil != err {
			log.Printf("iterate blockchain failed,err:%s", err.Error())
			return
		}

		log.Printf("Prev. hash: %x", block.PrevBlockHash)
		log.Printf("Data. %s", block.Data)
		log.Printf("Hash. %x", block.Hash)
		pow := NewProofOfWork(block)
		log.Printf("PoW: %s", strconv.FormatBool(pow.Validate()))

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) printUseage() string {
	return ""
}
