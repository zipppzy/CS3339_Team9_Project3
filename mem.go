package main

func MEM(ins Instruction) int {
	if ins.op == "LDUR" {
		return LoadMem(Registers[ins.rn] + int(ins.address)*4)
	} else if ins.op == "STUR" {
		StoreMem(Registers[ins.rn]+int(ins.address)*4, Registers[ins.rt])
	}
	return -1
}
