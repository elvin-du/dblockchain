package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		[]byte(data),
		prevBlockHash,
		[]byte{},
		0,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run() //开始挖矿

	block.Hash = hash
	block.Nonce = nonce

	block.SetHash()
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func DeserializeBlock(b []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)
	if nil != err {
		log.Printf("deserialize block failed,err:%s", err.Error())
		return nil, err
	}

	return &block, nil
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func (b *Block) Serialize() ([]byte, error) {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(b)
	if nil != err {
		log.Printf("serialize block failed,err:%s", err.Error())
		return nil, err
	}

	return buffer.Bytes(), nil
}
