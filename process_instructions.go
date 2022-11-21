package main

import (
	"fmt"
	"math"
	"strconv"
)

func ProcessInstructionList(list []Instruction) {
	breakHit := false
	for i := 0; i < len(list); i++ {
		if !breakHit {
			//use processing functions on instructions here
			translateToInt(&list[i])
			opcodeMasking(&list[i])
			opcodeTranslation(&list[i])
			switch list[i].instructionType {
			case "B":
				processBType(&list[i])
			case "I":
				processIType(&list[i])
			case "CB":
				processCBType(&list[i])
			case "IM":
				processIMType(&list[i])
			case "D":
				processDType(&list[i])
			case "R":
				processRType(&list[i])
			case "BREAK":
				breakHit = true
				BreakPoint = i
			}
		} else {
			list[i].instructionType = "MEM"
			var value uint64
			value, _ = strconv.ParseUint(list[i].rawInstruction, 2, 32)
			list[i].memValue = parse2Complement(value, 32)
			Mem[int(list[i].memLoc)] = parse2Complement(value, 32)
		}
	}
}

// translates raw instruction to an unsigned 64 bit int
func translateToInt(ins *Instruction) {
	i, err := strconv.ParseUint(ins.rawInstruction, 2, 64)
	if err == nil {
		ins.lineValue = i
	} else {
		fmt.Println(err)
	}
}

// 4292870144 mask for opcode(first 11 bits)
func opcodeMasking(ins *Instruction) {
	ins.opcode = (ins.lineValue & 4292870144) >> 21
}

// haven't handled BREAK yet
func opcodeTranslation(ins *Instruction) {
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

// \/\/ fill these out \/\/
func processRType(ins *Instruction) {
	//mask for bits 12 - 16
	ins.rm = uint8((ins.lineValue & 2031616) >> 16)
	//mask for bits 17 - 22
	ins.shamt = uint8((ins.lineValue & 64512) >> 10)
	//mask for bits 23 - 27
	ins.rn = uint8((ins.lineValue & 992) >> 5)
	//mask for bit 28 - 32
	ins.rd = uint8(ins.lineValue & 31)
}

func processIType(ins *Instruction) {
	//mask for bits 11 - 22
	ins.immediate = int16(parse2Complement((ins.lineValue&4193280)>>10, 12))
	//mask for bits 23 - 27
	ins.rn = uint8((ins.lineValue & 992) >> 5)
	//mask for bits 28 - 32
	ins.rd = uint8(ins.lineValue & 31)
}

func processCBType(ins *Instruction) {
	//mask for bits 9 - 27
	ins.offset = int32(parse2Complement((ins.lineValue&16777184)>>5, 19))
	//mask for bits 28 - 32
	ins.conditional = uint8(ins.lineValue & 31)
}

func processIMType(ins *Instruction) {
	//mask for bits 10 - 12
	ins.shiftCode = uint8((ins.lineValue & 6291456) >> 21)
	//mask for bits 13 - 27
	ins.field = uint32((ins.lineValue & 2097120) >> 5)
	//mask for bits 28 - 32
	ins.rd = uint8(ins.lineValue & 31)

}

func processDType(ins *Instruction) {
	//mask for bits 12 - 20
	ins.address = uint16((ins.lineValue & 2093056) >> 12)
	//mask for bits 21 - 22
	ins.op2 = uint8((ins.lineValue & 3072) >> 10)
	//mask for bits 23 - 27
	ins.rn = uint8((ins.lineValue & 992) >> 5)
	//mask for bit 28 - 32
	ins.rt = uint8(ins.lineValue & 31)
}

func processBType(ins *Instruction) {
	//mask for bits 7 - 32
	ins.offset = int32(parse2Complement(ins.lineValue&67108863, 26))
}

// parses 2's complement binary to an integer
func parse2Complement(i uint64, binaryLength uint) int64 {
	var out int64
	var xorValue int64
	out = int64(i)
	xorValue = (1 << binaryLength) - 1
	if (i >> (binaryLength - 1)) != 0 {
		out = ((out ^ xorValue) + 1) * -1
	}
	return out
}

// modifies registers and memory based on instruction given
func ExecuteInstruction(ins Instruction) {
	switch ins.op {
	case "B":
		branchTarget := PCindex + int(ins.offset)
		if branchTarget >= BreakPoint || branchTarget > 0 {
			PCindex = branchTarget
		}
	case "CBZ":
		branchTarget := PCindex + int(ins.offset)
		if Registers[ins.conditional] == 0 {
			if branchTarget >= BreakPoint || branchTarget > 0 {
				PCindex = branchTarget
			}
		}
	case "CBNZ":
		branchTarget := PCindex + int(ins.offset)
		if Registers[ins.conditional] != 0 {
			if branchTarget >= BreakPoint || branchTarget > 0 {
				PCindex = branchTarget
			}
		}
	case "ADDI":
		Registers[ins.rd] = Registers[ins.rn] + int(ins.immediate)
	case "SUBI":
		Registers[ins.rd] = Registers[ins.rn] - int(ins.immediate)
	case "AND":
		Registers[ins.rd] = Registers[ins.rn] & Registers[ins.rm]
	case "ADD":
		Registers[ins.rd] = Registers[ins.rn] + Registers[ins.rm]
	case "ORR":
		Registers[ins.rd] = Registers[ins.rn] | Registers[ins.rm]
	case "SUB":
		Registers[ins.rd] = Registers[ins.rn] - Registers[ins.rm]
	case "LSR":
		Registers[ins.rd] = Registers[ins.rn] >> Registers[ins.shamt]
	case "LSL":
		Registers[ins.rd] = Registers[ins.rn] << Registers[ins.shamt]
	case "ASR":
		Registers[ins.rd] = Registers[ins.rn] >> Registers[ins.shamt]
	case "EOR":
		Registers[ins.rd] = Registers[ins.rn] ^ Registers[ins.shamt]
	case "MOVZ":
		Registers[ins.rd] = int(ins.field) << 16 * int(ins.shiftCode)
	case "MOVK":
		masking := ^int(math.Pow(2, math.Ceil(math.Log2(float64(ins.field)))) - 1)
		value := int(ins.field) << 16 * int(ins.shiftCode)
		masking = masking << 16 * int(ins.shiftCode)
		Registers[ins.rd] = Registers[ins.rd] & masking
		Registers[ins.rd] = Registers[ins.rd] | value
	case "LDUR":
		Registers[ins.rt] = int(Mem[Registers[ins.rn]+int(ins.address)*4])
	case "STUR":
		Mem[Registers[ins.rn]+int(ins.address)*4] = int64(Registers[ins.rt])
	}
}
