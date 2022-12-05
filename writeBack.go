package main

func WB(ins Instruction, value int) {
	if ins.instructionType == "D" {
		Registers[ins.rt] = value
	} else {
		Registers[ins.rd] = value
	}
}
