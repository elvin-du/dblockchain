package main

import (
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

var (
	dbFile       = "blocks.db"
	blocksBucket = "blocksBucket"
)

type BlockChain struct {
	db            *bolt.DB
	lastBlockHash []byte
}

func NewBlockChain() (*BlockChain, error) {
	var lbh []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if nil != err {
		log.Printf("open db file(%s) failed,err:%s", dbFile, err.Error())
		return nil, err
	}
	//	defer db.Close()
	log.Printf("open db file(%s) success", dbFile)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if nil == b {
			log.Printf("bucket(%s) not found, so create a bucket and new genesis block", blocksBucket)
			genesis := NewGenesisBlock()
			log.Printf("new genesis block success")
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if nil != err {
				log.Printf("create bucket failed,err:%s", err.Error())
				return err
			}

			serializedBlock, err := genesis.Serialize()
			if nil != err {
				log.Printf("%+v block serialize failed,err:%s", genesis, err.Error())
				return err
			}
			err = b.Put(genesis.Hash, serializedBlock)
			if nil != err {
				log.Printf("put block into db failed,err:%s", err.Error())
				return err
			}
			log.Printf("put genesis block into chain")

			err = b.Put([]byte("l"), genesis.Hash)
			if nil != err {
				log.Printf("set last block hash failed,err:%s", err.Error())
				return err
			}
			lbh = genesis.Hash
		} else {
			lbh = b.Get([]byte("l"))
		}

		return nil
	})
	if nil != err {
		log.Printf("new blockchain failed,err:%s", err.Error())
		return nil, err
	}

	return &BlockChain{
		db:            db,
		lastBlockHash: lbh,
	}, nil
}

func (bc *BlockChain) AddBlock(data string) error {
	var lastBlockHash []byte

	bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastBlockHash = b.Get([]byte("l"))
		return nil
	})

	newBlock := NewBlock(data, lastBlockHash)

	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		serializedBlock, err := newBlock.Serialize()
		if nil != err {
			log.Printf("%+v block serialize failed,err:%s", newBlock, err.Error())
			return err
		}
		err = b.Put(newBlock.Hash, serializedBlock)
		if nil != err {
			log.Printf("put block into db failed,err:%s", err.Error())
			return err
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if nil != err {
			log.Printf("set last block hash failed,err:%s", err.Error())
			return err
		}
		bc.lastBlockHash = newBlock.Hash
		return nil
	})
	if nil != err {
		log.Printf("add block failed,err:%s", err.Error())
		return err
	}

	log.Printf("add block(data:%v) success", data)
	return nil
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{
		currentHash: bc.lastBlockHash,
		db:          bc.db,
	}
}

type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

var (
	E_BlockNotFound = errors.New("Block Not Found")
)

func (bci *BlockChainIterator) Next() (*Block, error) {
	var block *Block
	var err error = nil
	err = bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(bci.currentHash)
		if nil == encodedBlock {
			return E_BlockNotFound
		}

		block, err = DeserializeBlock(encodedBlock)
		if nil != err {
			log.Printf("deserialized block failed,err:%s", err.Error())
			return err
		}
		return nil
	})
	if nil != err {
		log.Printf("read block from db failed,err:%s", err.Error())
		return nil, err
	}

	bci.currentHash = block.PrevBlockHash

	return block, nil
}
