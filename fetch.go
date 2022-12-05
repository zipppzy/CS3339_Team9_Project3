package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Read line of binary code
func fetch(ins *Instruction) {
	readInstruct, err := os.Open("data.txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readInstruct)

	fileScanner.Split(bufio.ScanLines)
	ins.rawInstruction = fileScanner.Text()
}

// Translate to int
func decode(ins *Instruction) {
	i, err := strconv.ParseUint(ins.rawInstruction, 2, 64)
	if err == nil {
		ins.lineValue = i
	} else {
		ins.opcode = (ins.lineValue & 4292870144) >> 21
		if ins.opcode >= 160 && ins.opcode <= 191 {
			ins.op = "B"
			ins.instructionType = "B"
		} else if ins.opcode == 1104 {
			ins.op = "AND"
			ins.instructionType = "R"
		} else if ins.opcode == 1112 {
			ins.op = "ADD"
			ins.instructionType = "R"
		} else if ins.opcode >= 1160 && ins.opcode <= 1161 {
			ins.op = "ADDI"
			ins.instructionType = "I"
		} else if ins.opcode == 1360 {
			ins.op = "ORR"
			ins.instructionType = "R"
		} else if ins.opcode >= 1440 && ins.opcode <= 1447 {
			ins.op = "CBZ"
			ins.instructionType = "CB"
		} else if ins.opcode >= 1448 && ins.opcode <= 1455 {
			ins.op = "CBNZ"
			ins.instructionType = "CB"
		} else if ins.opcode == 1624 {
			ins.op = "SUB"
			ins.instructionType = "R"
		} else if ins.opcode >= 1672 && ins.opcode <= 1673 {
			ins.op = "SUBI"
			ins.instructionType = "I"
		} else if ins.opcode >= 1684 && ins.opcode <= 1687 {
			ins.op = "MOVZ"
			ins.instructionType = "IM"
		} else if ins.opcode >= 1940 && ins.opcode <= 1943 {
			ins.op = "MOVK"
			ins.instructionType = "IM"
		} else if ins.opcode == 1690 {
			ins.op = "LSR"
			ins.instructionType = "R"
		} else if ins.opcode == 1691 {
			ins.op = "LSL"
			ins.instructionType = "R"
		} else if ins.opcode == 1984 {
			ins.op = "STUR"
			ins.instructionType = "D"
		} else if ins.opcode == 1986 {
			ins.op = "LDUR"
			ins.instructionType = "D"
		} else if ins.opcode == 1692 {
			ins.op = "ASR"
			ins.instructionType = "R"
		} else if ins.opcode == 0 {
			ins.op = "NOP"
			ins.instructionType = "NOP"
		} else if ins.opcode == 1872 {
			ins.op = "EOR"
			ins.instructionType = "R"
		} else if ins.opcode == 2038 {
			ins.op = "BREAK"
			ins.instructionType = "BREAK"
		} else if ins.opcode == 0 {
			ins.op = "NOP"
			ins.instructionType = "NOP"
		} else {
			fmt.Println("Invalid opcode")
		}
	}
}
