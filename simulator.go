package main

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

}
