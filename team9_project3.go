package main

import (
	"flag"
)

type Instruction struct {
	rawInstruction  string
	lineValue       uint64
	memLoc          uint64
	opcode          uint64
	op              string
	instructionType string
	rm              uint8
	shamt           uint8
	rn              uint8
	rd              uint8
	rt              uint8
	op2             uint8
	address         uint16
	immediate       int16
	offset          int32
	conditional     uint8
	shiftCode       uint8
	field           uint32
	memValue        int64
}

var InstructionList []Instruction

// holds registers R0 - R31 (default 0)
var Registers [32]int

var Mem = make(map[int]int)

func main() {
	inputFilePathPtr := flag.String("i", "executionTest.txt", "input file path")
	outputFilePathPtr := flag.String("o", "out", "output file path")

	flag.Parse()
	//Inputs Command-Line
	ReadBinary(*inputFilePathPtr)

	ProcessInstructionList(InstructionList)

	WriteInstructions(*outputFilePathPtr+"_dis.txt", InstructionList)

	WriteInstructionExecution(*outputFilePathPtr+"_sim.txt", InstructionList)

}
