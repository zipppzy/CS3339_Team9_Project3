package main

var ALUResult uint16

func ALUOpI(entry1 uint16, immediate uint16) {
	switch opcode {
	case "ADDI":
		ALUResult = entry1 + immediate
	case "SUBI":
		ALUResult = entry1 - immediate
	}
}

func ALUOp(entry1 uint16, entry2 uint16) {
	switch opcode {
	case "AND":
		ALUResult = entry1 & entry2
	case "ADD":
		ALUResult = entry1 + entry2
	case "ORR":
		ALUResult = entry1 | entry2
	case "SUB":
		ALUResult = entry1 - entry2
	case "LSR":
		ALUResult = entry1 >> entry2
	case "LSL":
		ALUResult = entry1 << entry2
	case "ASR":
		ALUResult = entry1 >> entry2
	case "EOR":
		ALUResult = entry1 ^ entry2
	}
}
