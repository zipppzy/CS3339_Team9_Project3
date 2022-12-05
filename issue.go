package main

// moves up to two instructions into either the PreAluBuffer or the PreMemBuffer
func Issue() {
	moved := 0
	buff := make([]int, 0)
	for len(PreIssueBuff) > 0 {
		buff = append(buff, <-PreIssueBuff)
	}
	for index, element := range buff {
		if moved >= 2 {
			break
		}
		structuralHazard := false
		if InstructionList[element].instructionType == "D" && len(PreMemBuff) >= 2 {
			structuralHazard = true
		}

		if InstructionList[element].instructionType != "D" && len(PreALUBuff) >= 2 {
			structuralHazard = true
		}

		if !structuralHazard {
			rawHazard := false
			//check preIssueBuff before index for RAW hazard
			for i := index - 1; i >= 0; i-- {
				if isRawHazard(InstructionList[i], InstructionList[element]) {
					rawHazard = true
					break
				}
			}

			//check PreMemBuff for RAW hazard
			if !rawHazard {
				//make temp slice to store PreMemBuff
				tempBuff := make([]int, 0)
				for len(PreMemBuff) > 0 {
					tempBuff = append(tempBuff, <-PreMemBuff)
				}

				for _, e := range tempBuff {
					if isRawHazard(InstructionList[e], InstructionList[element]) {
						rawHazard = true
						break
					}
				}
				//return tempBuff to PreMemBuff
				for _, e := range tempBuff {
					PreMemBuff <- e
				}
			}

			//check PreALUBuff for Raw hazard
			if !rawHazard {
				//make temp slice to store PreALUBuff
				tempBuff := make([]int, 0)
				for len(PreALUBuff) > 0 {
					tempBuff = append(tempBuff, <-PreALUBuff)
				}

				for _, e := range tempBuff {
					if isRawHazard(InstructionList[e], InstructionList[element]) {
						rawHazard = true
						break
					}
				}
				//return tempBuff to PreALUBuff
				for _, e := range tempBuff {
					PreALUBuff <- e
				}
			}

			//check PostALUBuff for Raw hazard
			if !rawHazard && len(postALUBuff) != 0 {
				temp := <-postALUBuff
				if isRawHazard(InstructionList[temp[0]], InstructionList[element]) {
					rawHazard = true
				}
				postALUBuff <- temp
			}

			//if there is no hazard write to appropriate buff
			if !rawHazard {
				if InstructionList[element].instructionType == "D" {
					PreMemBuff <- element
					buff[index] = -1
					moved++
				} else {
					PreALUBuff <- element
					buff[index] = -1
					moved++
				}
			}
		}
	}

	//return remaining buff to PreIssueBuff
	for _, e := range buff {
		if e != -1 {
			PreMemBuff <- e
		}
	}
}

func isRawHazard(ins1 Instruction, ins2 Instruction) bool {
	//destination register for ins1
	var dest int = -1
	//operand 1 for ins2
	var op1 int = -1
	//operand 2 for ins2
	var op2 int = -1

	if ins1.op == "LDUR" {
		dest = int(ins1.rt)
	} else {
		dest = int(ins1.rd)
	}

	if ins2.instructionType == "D" || ins2.instructionType == "I" {
		op1 = int(ins2.rn)
		op2 = int(ins2.rn)
	} else {
		op1 = int(ins2.rm)
		op2 = int(ins2.rn)
	}

	return dest == op1 || dest == op2
}
