package main

//var preIssueBuff = [4]int{-1, -1, -1, -1}
//var preMemBuff = [2]int{-1, -1}
//var preALUBuff = [2]int{-1, -1}
//var postMemBuff = [2]int{-1, -1}
//var postALUBuff = [2]int{-1, -1}

var preIssueBuff = make(chan int, 4)
var preMemBuff = make(chan int, 2)
var preALUBuff = make(chan int, 2)

// stores 1 array with instruction index at [0] and value at [1]
var postMemBuff = make(chan [2]int, 1)
var postALUBuff = make(chan [2]int, 1)

func Simulate() {

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

	if len(preALUBuff) != 0 {
		insIndex := <-preALUBuff
		var AluOut = [2]int{insIndex, ALU(InstructionList[insIndex])}
		postALUBuff <- AluOut
	}

	if len(preMemBuff) != 0 {
		insIndex := <-preMemBuff
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
			if len(preMemBuff) == 0 {
				preMemBuff <- insIndex
			} else if len(preMemBuff) == 1 {
				tempInt := <-preMemBuff
				preMemBuff <- insIndex
				preMemBuff <- tempInt
			}
		}
	}

}
