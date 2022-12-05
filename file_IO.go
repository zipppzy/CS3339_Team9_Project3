package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// index in instruction list
var PCindex = 0

// index where Break instruction is
var BreakPoint int

// ReadBinary reads text file and makes Instructions and adds them to the InstructionList
func ReadBinary(filePath string) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var pc uint64
	pc = 96
	for scanner.Scan() {
		InstructionList = append(InstructionList, Instruction{rawInstruction: scanner.Text(), memLoc: pc})
		pc += 4
	}
}

func WriteInstructions(filePath string, list []Instruction) {
	f, err := os.Create(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for i := 0; i < len(list); i++ {
		switch list[i].instructionType {
		case "B":
			//write binary with spaces
			_, err := fmt.Fprintf(f, "%s %s\t", list[i].rawInstruction[0:6], list[i].rawInstruction[6:32])
			//write memLoc and opcode
			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].memLoc, list[i].op)
			//write operands
			_, err = fmt.Fprintf(f, "#%d\n", list[i].offset)
			if err != nil {
				log.Fatal(err)
			}
		case "I":
			//write binary with spaces
			_, err := fmt.Fprintf(f, "%s %s %s %s\t", list[i].rawInstruction[0:10], list[i].rawInstruction[10:22], list[i].rawInstruction[22:27], list[i].rawInstruction[27:32])
			//write memLoc and opcode
			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].memLoc, list[i].op)
			//write operands
			_, err = fmt.Fprintf(f, "R%d, R%d, #%d\n", list[i].rd, list[i].rn, list[i].immediate)
			if err != nil {
				log.Fatal(err)
			}

		case "CB":
			//write binary with spaces
			_, err := fmt.Fprintf(f, "%s %s %s\t", list[i].rawInstruction[0:8], list[i].rawInstruction[8:27], list[i].rawInstruction[27:32])
			//write memLoc and opcode
			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].memLoc, list[i].op)
			//write operands
			_, err = fmt.Fprintf(f, "R%d, #%d\n", list[i].conditional, list[i].offset)
			if err != nil {
				log.Fatal(err)
			}
		case "IM":
			//write binary with spaces
			_, err := fmt.Fprintf(f, "%s %s %s %s\t", list[i].rawInstruction[0:9], list[i].rawInstruction[9:12], list[i].rawInstruction[12:27], list[i].rawInstruction[27:32])
			//write memLoc and opcode
			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].memLoc, list[i].op)
			//write operands
			_, err = fmt.Fprintf(f, "R%d, %d, LSL %d\n", list[i].rd, list[i].field, list[i].shiftCode)
			if err != nil {
				log.Fatal(err)
			}
			// I am not sure about D too
		case "D":
			//write binary with spaces
			_, err := fmt.Fprintf(f, "%s %s %s %s %s\t", list[i].rawInstruction[0:11], list[i].rawInstruction[11:20], list[i].rawInstruction[20:22], list[i].rawInstruction[22:27], list[i].rawInstruction[27:32])
			//write memLoc and opcode
			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].memLoc, list[i].op)
			//write operands
			_, err = fmt.Fprintf(f, "R%d, [R%d, #%d]\n", list[i].rt, list[i].rn, list[i].address)
			if err != nil {
				log.Fatal(err)
			}
		case "R":
			//write binary with spaces
			_, err := fmt.Fprintf(f, "%s %s %s %s %s\t", list[i].rawInstruction[0:11], list[i].rawInstruction[11:16], list[i].rawInstruction[16:22], list[i].rawInstruction[22:27], list[i].rawInstruction[27:32])
			//write memLoc and opcode
			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].memLoc, list[i].op)
			//write operands
			_, err = fmt.Fprintf(f, "R%d, R%d, ", list[i].rd, list[i].rn)
			if list[i].op == "LSL" || list[i].op == "ASR" || list[i].op == "LSR" {
				_, err = fmt.Fprintf(f, "#%d\n", list[i].shamt)
			} else {
				_, err = fmt.Fprintf(f, "R%d\n", list[i].rm)
			}
			if err != nil {
				log.Fatal(err)
			}
		case "BREAK":
			_, err := fmt.Fprintf(f, "%s\t%d\tBREAK\n", list[i].rawInstruction, list[i].memLoc)
			if err != nil {
				log.Fatal(err)
			}
		case "MEM":
			_, err := fmt.Fprintf(f, "%s\t%d\t%d\n", list[i].rawInstruction, list[i].memLoc, list[i].memValue)
			if err != nil {
				log.Fatal(err)
			}
		case "NOP":
			_, err := fmt.Fprintf(f, "%s\t%d\t%s\n", list[i].rawInstruction, list[i].memLoc, list[i].op)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// executes all instructions and writes state of registers and memory at each step
func WriteInstructionExecution(filePath string, list []Instruction) {

	f, err := os.Create(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	var cycle = 1

	for PCindex <= BreakPoint {
		_, err := fmt.Fprintf(f, "\n====================\n")
		// print cycle and memory location of instruction and op
		_, err = fmt.Fprintf(f, "cycle:%d\t%d\t%s\t", cycle, list[PCindex].memLoc, list[PCindex].op)

		//prints just the operands
		switch list[PCindex].instructionType {
		case "B":
			//write operands
			_, err = fmt.Fprintf(f, "#%d\n", list[PCindex].offset)
			if err != nil {
				log.Fatal(err)
			}
		case "I":
			//write operands
			_, err = fmt.Fprintf(f, "R%d, R%d, #%d\n", list[PCindex].rd, list[PCindex].rn, list[PCindex].immediate)
			if err != nil {
				log.Fatal(err)
			}

		case "CB":
			//write operands
			_, err = fmt.Fprintf(f, "R%d, #%d\n", list[PCindex].conditional, list[PCindex].offset)
			if err != nil {
				log.Fatal(err)
			}
		case "IM":
			//write operands
			_, err = fmt.Fprintf(f, "R%d, %d, LSL %d\n", list[PCindex].rd, list[PCindex].field, list[PCindex].shiftCode)
			if err != nil {
				log.Fatal(err)
			}
		case "D":
			//write operands
			_, err = fmt.Fprintf(f, "R%d, [R%d, #%d]\n", list[PCindex].rt, list[PCindex].rn, list[PCindex].address)
			if err != nil {
				log.Fatal(err)
			}
		case "R":
			//write operands
			_, err = fmt.Fprintf(f, "R%d, R%d, ", list[PCindex].rd, list[PCindex].rn)
			if list[PCindex].op == "LSL" || list[PCindex].op == "ASR" || list[PCindex].op == "LSR" {
				_, err = fmt.Fprintf(f, "#%d\n", list[PCindex].shamt)
			} else {
				_, err = fmt.Fprintf(f, "R%d\n", list[PCindex].rm)
			}
			if err != nil {
				log.Fatal(err)
			}
		default:
			_, err = fmt.Fprintf(f, "\n")
			if err != nil {
				log.Fatal(err)
			}
		}

		ExecuteInstruction(list[PCindex])

		_, err = fmt.Fprintf(f, "\nregisters:\n")

		//prints registers
		_, err = fmt.Fprintf(f, "r00:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", Registers[0], Registers[1], Registers[2], Registers[3], Registers[4], Registers[5], Registers[6], Registers[7])
		_, err = fmt.Fprintf(f, "r08:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", Registers[8], Registers[9], Registers[10], Registers[11], Registers[12], Registers[13], Registers[14], Registers[15])
		_, err = fmt.Fprintf(f, "r16:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", Registers[16], Registers[17], Registers[18], Registers[19], Registers[20], Registers[21], Registers[22], Registers[23])
		_, err = fmt.Fprintf(f, "r24:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n\n", Registers[24], Registers[25], Registers[26], Registers[27], Registers[28], Registers[29], Registers[30], Registers[31])

		_, err = fmt.Fprintf(f, "data:")
		key := int(list[BreakPoint+1].memLoc)
		// find the largest memory location with something assigned
		largestKey := key
		for i := range Mem {
			if i > largestKey {
				largestKey = i
			}
		}

		for key <= largestKey {
			if (key-(int(list[BreakPoint].memLoc)+4))%32 == 0 {
				_, err = fmt.Fprintf(f, "\n%d:\t", key)
			}
			_, err = fmt.Fprintf(f, "%d\t", Mem[key])
			key += 4
		}
		PCindex++
		cycle++
	}
}
