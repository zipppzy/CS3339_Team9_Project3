package main

type block struct {
	valid int
	dirty int
	tag   int
	word1 int
	word2 int
	//value of block as a 64 bit num
	value int
}

var CacheSets [4][2]block

// indicates which block in each set was least recently used
var lruBits = [4]int{0, 0, 0, 0}

var tagMask = 4294967264
var setMask = 24
var word1Mask = 4294967295
var word2Mask = 4294967295 << 32

//Set# = (address/4)%4

func StoreMem(address int, value int) {
	cacheHit, blockNum := checkCacheHit(address)
	var setNum = (address & setMask) >> 3
	var word1Val = value & word1Mask >> 32
	var word2Val = value & word2Mask
	if cacheHit {
		CacheSets[setNum][blockNum].word1 = word1Val
		CacheSets[setNum][blockNum].word2 = word2Val
		CacheSets[setNum][blockNum].value = value
	} else {
		var tag = (address & tagMask) >> 5
		CacheSets[setNum][lruBits[setNum]] = block{valid: 1, tag: tag, word1: word1Val, word2: word2Val, value: value}
		//flip lruBit
		if lruBits[setNum] == 0 {
			lruBits[setNum] = 1
		} else {
			lruBits[setNum] = 0
		}
	}
	Mem[address] = value
}

func LoadMem(address int) int {
	var setNum = (address & setMask) >> 3
	var tag = (address & tagMask) >> 5
	cacheHit, blockNum := checkCacheHit(address)
	if cacheHit {
		return CacheSets[setNum][blockNum].value
	} else {
		//find which two addresses to load to cache
		var address1 int
		var address2 int
		if address%8 == 0 {
			address1 = address
			address2 = address + 4
		} else {
			address1 = address - 4
			address2 = address
		}

		combinedVal := (Mem[address1] << 32) + Mem[address2]

		//load into cache
		CacheSets[setNum][lruBits[setNum]] = block{valid: 1, tag: tag, word1: Mem[address1], word2: Mem[address2], value: combinedVal}

		//flip lruBit
		if lruBits[setNum] == 0 {
			lruBits[setNum] = 1
		} else {
			lruBits[setNum] = 0
		}

		return combinedVal
	}
}

func checkCacheHit(address int) (bool, int) {
	var tag = (address & tagMask) >> 5
	var setNum = (address & setMask) >> 3

	//checks if tag is equal and if valid bit is 1 for both blocks
	if (CacheSets[setNum][0].tag == tag) && (CacheSets[setNum][0].valid == 1) {
		return true, 0
	} else if (CacheSets[setNum][1].tag == tag) && (CacheSets[setNum][0].valid == 1) {
		return true, 1
	} else {
		return false, -1
	}
}
