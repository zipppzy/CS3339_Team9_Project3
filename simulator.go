package main

type Block struct {
	valid int
	dirty int
	tag   int
	word1 int
	word2 int
}

var CacheSets [4][2]Block
