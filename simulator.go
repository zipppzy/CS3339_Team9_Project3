package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

//var preIssueBuff = [4]int{-1, -1, -1, -1}
//var preMemBuff = [2]int{-1, -1}
//var preALUBuff = [2]int{-1, -1}
//var postMemBuff = [2]int{-1, -1}
//var postALUBuff = [2]int{-1, -1}

var PreIssueBuff = make(chan int, 4)
var PreMemBuff = make(chan int, 2)
var PreALUBuff = make(chan int, 2)

// stores 1 array with instruction index at [0] and value at [1]
var postMemBuff = make(chan [2]int, 1)
var postALUBuff = make(chan [2]int, 1)

var Break bool = false
var cycleNum = 0

func Simulate(filePath string) {
	f, err := os.Create(filePath)
	defer f.Close()
	Cycle()
	PrintCycle(f)
	Cycle()
	PrintCycle(f)
	for len(postALUBuff) != 0 || len(postMemBuff) != 0 || len(PreALUBuff) != 0 || len(PreMemBuff) != 0 || len(PreIssueBuff) != 0 {
		Cycle()
		PrintCycle(f)
		cycleNum++
	}
	if err != nil {
		log.Fatal(err)
	}
}

func Cycle() {
	//WB both postALUBuff and postMemBuff
	if len(postALUBuff) != 0 {
		buff := <-postALUBuff
		WB(InstructionList[buff[1]], buff[0])
	}
	if len(postMemBuff) != 0 {
		buff := <-postMemBuff
		WB(InstructionList[buff[1]], buff[0])
	}

	if len(PreALUBuff) != 0 {
		insIndex := <-PreALUBuff
		var AluOut = [2]int{insIndex, ALU(InstructionList[insIndex])}
		postALUBuff <- AluOut
	}

	if len(PreMemBuff) != 0 {
		insIndex := <-PreMemBuff
		//check for cache hit
		cacheHit, _ := CheckCacheHit(Registers[InstructionList[insIndex].rn] + int(InstructionList[insIndex].address)*4)
		if cacheHit {
			//proceed as expected
			var MemOut = [2]int{insIndex, MEM(InstructionList[insIndex])}
			if MemOut[1] != -1 {
				postMemBuff <- MemOut
			}
		} else {
			//change MEM and cache
			MEM(InstructionList[insIndex])
			//put insIndex back in preMemBuff queue in correct order
			if len(PreMemBuff) == 0 {
				PreMemBuff <- insIndex
			} else if len(PreMemBuff) == 1 {
				tempInt := <-PreMemBuff
				PreMemBuff <- insIndex
				PreMemBuff <- tempInt
			}
		}
	}

	if len(PreIssueBuff) != 0 {
		Issue()
	}

	if !Break {
		for i := 0; i < 2; i++ {
			if Fetch() {
				break
			}
		}
	}
}

func PrintCycle(f *os.File) {
	_, err := fmt.Fprintf(f, "--------------------\n")
	_, err = fmt.Fprintf(f, "Cycle: %d\n", cycleNum)
	_, err = fmt.Fprintf(f, "Pre-Issue Buffer:\n")

	for y := 0; y < 4; y++ {
		_, err = fmt.Fprintf(f, "\tEntry %d:\t[instruction]\n", y)
	}
	_, err = fmt.Fprintf(f, "Pre_ALU Queue:\n")
	for y := 0; y < 2; y++ {
		_, err = fmt.Fprintf(f, "\tEntry %d:\t[instruction]\n", y)
	}
	_, err = fmt.Fprintf(f, "Post_ALU Queue:\n")
	_, err = fmt.Fprintf(f, "\tEntry 0:\t[instruction]\n")
	_, err = fmt.Fprintf(f, "Pre_MEM Queue:\n")
	for y := 0; y < 2; y++ {
		_, err = fmt.Fprintf(f, "\tEntry %d:\t[instruction]\n", y)
	}
	_, err = fmt.Fprintf(f, "Post_MEM Queue:\n")
	_, err = fmt.Fprintf(f, "\tEntry 0:\t[instruction]\n")
	_, err = fmt.Fprintf(f, "Registers\n")
	_, err = fmt.Fprintf(f, "r00:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", Registers[0], Registers[1], Registers[2], Registers[3], Registers[4], Registers[5], Registers[6], Registers[7])
	_, err = fmt.Fprintf(f, "r08:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", Registers[8], Registers[9], Registers[10], Registers[11], Registers[12], Registers[13], Registers[14], Registers[15])
	_, err = fmt.Fprintf(f, "r16:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", Registers[16], Registers[17], Registers[18], Registers[19], Registers[20], Registers[21], Registers[22], Registers[23])
	_, err = fmt.Fprintf(f, "r24:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n\n", Registers[24], Registers[25], Registers[26], Registers[27], Registers[28], Registers[29], Registers[30], Registers[31])
	_, err = fmt.Fprintf(f, "Cache\n")
	for i := 0; i < 4; i++ {
		_, err = fmt.Fprintf(f, "Set %d: LRU=%d\n", i, LruBits[i])
		for j := 0; j < 2; j++ {
			_, err = fmt.Fprintf(f, "\tEntry %d:\t[(%d, %d, %d)<%s,%s>]\n", j, CacheSets[i][j].valid, CacheSets[i][j].dirty, CacheSets[i][j].tag, strconv.FormatInt(int64(CacheSets[i][j].word1), 2), strconv.FormatInt(int64(CacheSets[i][j].word2), 2))
		}
	}

	_, err = fmt.Fprintf(f, "\nData")
	key := int(InstructionList[BreakPoint+1].memLoc)
	largestKey := key
	for i := range Mem {
		if i > largestKey {
			largestKey = i
		}
	}

	for key <= largestKey {
		if (key-(int(InstructionList[BreakPoint].memLoc)+4))%32 == 0 {
			_, err = fmt.Fprintf(f, "\n%d:\t", key)
		}
		_, err = fmt.Fprintf(f, "%d\t", Mem[key])
		key += 4
	}
	if err != nil {
		log.Fatal(err)
	}
}
