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
var LruBits = [4]int{0, 0, 0, 0}

var tagMask = 4294967264
var setMask = 24

// var word1Mask = 4294967295
var word2Mask = 4294967295

//Set# = (address/4)%4

func StoreMem(address int, value int) {
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
		CacheSets[setNum][LruBits[setNum]] = block{valid: 1, dirty: 1, tag: tag, word1: word1Val, word2: word2Val, value: value}
		//flip lruBit
		if LruBits[setNum] == 0 {
			LruBits[setNum] = 1
		} else {
			LruBits[setNum] = 0
		}
	}
	Mem[address] = value
}

func LoadMem(address int) int {
	var setNum = (address & setMask) >> 3
	var tag = (address & tagMask) >> 5
	cacheHit, blockNum := CheckCacheHit(address)
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
		CacheSets[setNum][LruBits[setNum]] = block{valid: 1, tag: tag, word1: Mem[address1], word2: Mem[address2], value: combinedVal}

		//flip lruBit
		if LruBits[setNum] == 0 {
			LruBits[setNum] = 1
		} else {
			LruBits[setNum] = 0
		}

		return combinedVal
	}
}

func CheckCacheHit(address int) (bool, int) {
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

func LoadInstruction(ins Instruction) {
	var setNum = (int(ins.memLoc) & setMask) >> 3
	var tag = (int(ins.memLoc) & tagMask) >> 5

	var word1Val = int(ins.lineValue) >> 32
	var word2Val = int(ins.lineValue) & word2Mask

	CacheSets[setNum][LruBits[setNum]] = block{valid: 1, tag: tag, word1: word1Val, word2: word2Val, value: int(ins.lineValue)}
	//flip lruBit
	if LruBits[setNum] == 0 {
		LruBits[setNum] = 1
	} else {
		LruBits[setNum] = 0
	}
}
