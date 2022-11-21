package main

type Block struct {
	valid int
	dirty int
	tag   int
	word1 int
	word2 int
}

var CacheSets [4][2]Block

// indicates which block in each se was least recently used
var lruBits = [4]int{0, 0, 0, 0}

//Set# = (address/4)%4

func WriteMem(address int, value int) {

}

func LoadMem(address int) {
	//find which two addresses to load to cache
	var address1 int
	var address2 int
	if address%8 == 0 {
		address1 = address
		address2 = address + 4
	} else {
		adress1 = address - 4
		adress2 = address
	}

}

func checkCacheHit(address int) {

}
