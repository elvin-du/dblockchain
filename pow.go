package main

import (
	"bytes"
	"crypto/sha256"
	"log"
	"math"
	"math/big"
)

const (
	targetBits = 24
)

type ProofOfWork struct {
	block  *Block //要验证的区块
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	//把targetBits位数转变为可以比较大小的数值,是目标值的上界
	target := big.NewInt(1)
	//因为用的是sha256，所以用256减去targetBits
	target.Lsh(target, uint(256-targetBits))

	return &ProofOfWork{
		b,
		target,
	}
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

//算出nonce和hash值
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int //hash的整型表示
	var hash [32]byte
	nonce := 0
	maxNonce := math.MaxInt64

	log.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:]) //把hash转变为整型

		if hashInt.Cmp(pow.target) == -1 {
			log.Printf("hash. %x", hash)
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
