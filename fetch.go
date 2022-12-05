package main

var pc int = 0

func Fetch() bool {
	ins := InstructionList[pc]
	cacheHit, _ := CheckCacheHit(int(ins.memLoc))
	branched := false
	if cacheHit {
		if ins.op == "B" {
			pc = pc + int(ins.offset)
			branched = true
		} else if ins.op == "CBZ" {
			if Registers[ins.conditional] == 0 {
				pc = pc + int(ins.offset)
				branched = true
			}
		} else if ins.op == "CBNZ" {
			if Registers[ins.conditional] != 0 {
				pc = pc + int(ins.offset)
				branched = true
			}
		} else if ins.op == "NOP" {

		} else if ins.op == "BREAK" {
			Break = true
		} else {
			PreIssueBuff <- pc
			pc++
		}
	} else {
		LoadInstruction(ins)
	}
	return branched
}
