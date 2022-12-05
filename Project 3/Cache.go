package main

type Block struct {
	valid int
	dirty int
	tag   int
	word1 int
	word2 int
}

var CacheSets [4][2]Block
var LRUbits = [4]int{0, 0, 0, 0}

func WriteMem(address int, value int) {
	
}

func LoadMem(address int) {

}

func checkCacheHit(address int) {

}
