package main

var ALUResult int

func ALU(ins Instruction) int {
	if ins.instructionType == "I" {
		ALUOpI(Registers[ins.rn], int(ins.immediate), ins.op)
	} else {
		ALUOp(Registers[ins.rn], Registers[ins.rm], ins.op)
	}
	return ALUResult
}
func ALUOpI(entry1 int, immediate int, op string) {
	switch op {
	case "ADDI":
		ALUResult = entry1 + immediate
	case "SUBI":
		ALUResult = entry1 - immediate
	}
}

func ALUOp(entry1 int, entry2 int, op string) {
	switch op {
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
