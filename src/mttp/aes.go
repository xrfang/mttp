package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"strconv"
	"time"
)

func genKey() []byte {
	pass := cf.PASSWORD
	if pass == "" {
		pass = strconv.Itoa(int(time.Now().UnixNano()))
	}
	pass += verinfo()
	h := sha256.New()
	h.Write([]byte(pass))
	return h.Sum(nil)
}

func Encrypt(plain []byte) []byte {
	key := genKey()
	block, err := aes.NewCipher(key)
	assert(err)
	bsize := block.BlockSize()
	padlen := bsize - len(plain)%bsize
	padtxt := bytes.Repeat([]byte{byte(padlen)}, padlen)
	pdata := append(plain, padtxt...)
	cbc := cipher.NewCBCEncrypter(block, key[:bsize])
	crypted := make([]byte, len(pdata))
	cbc.CryptBlocks(crypted, pdata)
	return crypted
}

func Decrypt(crypted []byte) []byte {
	key := genKey()
	block, err := aes.NewCipher(key)
	assert(err)
	bsize := block.BlockSize()
	cbc := cipher.NewCBCDecrypter(block, key[:bsize])
	plain := make([]byte, len(crypted))
	cbc.CryptBlocks(plain, crypted)
	L := len(plain)
	x := int(plain[L-1])
	return plain[:(L - x)]
}
