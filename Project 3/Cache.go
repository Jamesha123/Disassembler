package main

var CacheSets [4][2]Block
var LRUBits = [4]int{0, 0, 0, 0}
var tagMask = 4294967264
var setMask = 24
var word2Mask = 4294967295

// Load from cache
func LDURMem(address int) int {
	var tag = (address & tagMask) >> 5
	var setNum = (address & setMask) >> 3
	cacheHit, blockNum := CheckCacheHit(address)
	if cacheHit {
		return CacheSets[setNum][blockNum].value
	} else {
		var address1 int
		var address2 int
		if address%8 == 0 { // is alligned
			address1 = address
			address2 = address + 4
		} else { // isnt alligned
			address1 = address - 4
			address2 = address
		}

		concatValues := (MemoryIndex[address1] << 32) + MemoryIndex[address2]
		CacheSets[setNum][LRUBits[setNum]] = Block{valid: 1,
			tag:   tag,
			word1: MemoryIndex[address1],
			word2: MemoryIndex[address2],
			value: concatValues,
		}

		if LRUBits[setNum] == 0 {
			LRUBits[setNum] = 1
		} else {
			LRUBits[setNum] = 0
		}

		return concatValues
	}
}

// stores into cache
func STURMem(address int, value int) {
	cacheHit, blockNum := CheckCacheHit(address)
	var setNum = (address & setMask) >> 3
	var word1Val = value >> 32
	var word2Val = value & word2Mask
	if cacheHit {
		CacheSets[setNum][blockNum].word1 = word1Val
		CacheSets[setNum][blockNum].word2 = word2Val
		CacheSets[setNum][blockNum].value = value
		CacheSets[setNum][blockNum].dirty = 1
	} else {
		var tag = (address & tagMask) >> 5
		CacheSets[setNum][LRUBits[setNum]] = Block{valid: 1,
			dirty: 1,
			tag:   tag,
			word1: word1Val,
			word2: word2Val,
			value: value,
		}
		if LRUBits[setNum] == 0 {
			LRUBits[setNum] = 1
		} else {
			LRUBits[setNum] = 0
		}
	}
	MemoryIndex[address] = value
}

func CheckCacheHit(address int) (bool, int) {
	var tag = (address & tagMask) >> 5
	var setNum = (address & setMask) >> 3

	//checks if tag is valid bit
	if (CacheSets[setNum][0].tag == tag) && (CacheSets[setNum][0].valid == 1) {
		return true, 0
	} else if (CacheSets[setNum][1].tag == tag) && (CacheSets[setNum][0].valid == 1) {
		return true, 1
	} else {
		return false, -1
	}
}
