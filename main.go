package main

import (
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	bc, err := NewBlockChain()
	if nil != err {
		log.Printf("new blockchain failed,err:%s", err.Error())
		return
	}
	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()
}
