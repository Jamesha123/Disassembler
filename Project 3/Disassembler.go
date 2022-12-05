package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

/*
*************************************
// functions
*************************************
*/

func processInput(list []Instruction) {

	for i := 0; i < len(list); i++ {
		if !Breaknow {
			BreakPoint = i
			convertToInt(&list[i])
			opcodeMasking(&list[i])
			opcodeMatching(&list[i])
			switch list[i].typeofInstruction {
			case "R":
				RTypeFormat(&list[i])
			case "D":
				DTypeFormat(&list[i])
			case "I":
				ITypeFormat(&list[i])
			case "B":
				BTypeFormat(&list[i])
			case "CB":
				CBTypeFormat(&list[i])
			case "IM":
				IMTypeFormat(&list[i])
			case "BREAK":
				Breaknow = true
				PC = list[i].programCnt + 4
				BreakPoint = i
			}
		} else {
			list[i].typeofInstruction = "NUM"
			var val uint64
			val, _ = strconv.ParseUint(list[i].rawInstruction, 2, 32)
			list[i].bitValue = int(TwoComplement(val, 32))

			counter := list[i].programCnt
			value := int(list[i].bitValue)
			temp := []int{counter, value}
			Data = append(Data, temp)
		}

		MemoryIndex[int(list[i].programCnt)] = list[i].bitValue
	}

}

// mask the opcode (11 bits)
func opcodeMasking(list *Instruction) {
	// mask the bits by 0xFFE00000 = 4292870144 to get first 11 bits
	list.opcode = (list.linevalue & 4292870144) >> 21
}

// converts opcode to int
func convertToInt(list *Instruction) {
	i, err := strconv.ParseUint(list.rawInstruction, 2, 64)
	//fmt.Println(i)
	if err == nil {
		list.linevalue = i
	} else {
		fmt.Println(err)
	}
}

// match the converted opcode to the arm instruction
func opcodeMatching(list *Instruction) {
	if list.opcode >= 160 && list.opcode <= 191 {
		list.op = "B"
		list.typeofInstruction = "B"
	} else if list.opcode == 1104 {
		list.op = "AND"
		list.typeofInstruction = "R"
	} else if list.opcode == 1112 {
		list.op = "ADD"
		list.typeofInstruction = "R"
	} else if list.opcode >= 1160 && list.opcode <= 1161 {
		list.op = "ADDI"
		list.typeofInstruction = "I"
	} else if list.opcode == 1360 {
		list.op = "ORR"
		list.typeofInstruction = "R"
	} else if list.opcode >= 1440 && list.opcode <= 1447 {
		list.op = "CBZ"
		list.typeofInstruction = "CB"
	} else if list.opcode >= 1448 && list.opcode <= 1455 {
		list.op = "CBNZ"
		list.typeofInstruction = "CB"
	} else if list.opcode == 1624 {
		list.op = "SUB"
		list.typeofInstruction = "R"
	} else if list.opcode >= 1672 && list.opcode <= 1673 {
		list.op = "SUBI"
		list.typeofInstruction = "I"
	} else if list.opcode >= 1684 && list.opcode <= 1687 {
		list.op = "MOVZ"
		list.typeofInstruction = "IM"
	} else if list.opcode >= 1940 && list.opcode <= 1943 {
		list.op = "MOVK"
		list.typeofInstruction = "IM"
	} else if list.opcode == 1690 {
		list.op = "LSR"
		list.typeofInstruction = "R"
	} else if list.opcode == 1691 {
		list.op = "LSL"
		list.typeofInstruction = "R"
	} else if list.opcode == 1984 {
		list.op = "STUR"
		list.typeofInstruction = "D"
	} else if list.opcode == 1986 {
		list.op = "LDUR"
		list.typeofInstruction = "D"
	} else if list.opcode == 1692 {
		list.op = "ASR"
		list.typeofInstruction = "R"
	} else if list.opcode == 0 && list.linevalue == 0 {
		list.op = "NOP"
		list.typeofInstruction = "NOP"
	} else if list.opcode == 1872 {
		list.op = "EOR"
		list.typeofInstruction = "R"
	} else if list.opcode == 2038 {
		list.op = "BREAK"
		list.typeofInstruction = "BREAK"
	} else {
		fmt.Println("Invalid Opcode")
	}
}

func readInstruction(filePath string) {

	fmt.Println("Opening a file")
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Could not open file: ", err)
	}

	// Close file after main runs
	defer file.Close()
	// program counter
	var pc = 96

	// Read in file with scanner
	//InputParsed := []Instruction{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ns := Instruction{rawInstruction: scanner.Text(), programCnt: pc}
		InputParsed = append(InputParsed, ns)
		pc += 4
	}
	if err := scanner.Err(); err != nil {
		fmt.Print(err)
	}
}

func writeInstruction(filePath string, list []Instruction) {
	f, err := os.Create(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for i := 0; i < len(list); i++ {
		switch list[i].typeofInstruction {
		case "R":
			// Prints out bits
			_, err := fmt.Fprintf(f, "%s %s %s %s %s\t", list[i].rawInstruction[0:11],
				list[i].rawInstruction[11:16], list[i].rawInstruction[16:22],
				list[i].rawInstruction[22:27], list[i].rawInstruction[27:32])
			if err != nil {
				log.Fatal(err)
			}
			// Prints out pc and opcode instruction
			_, err = fmt.Fprintf(f, "%d\t%s ", list[i].programCnt, list[i].op)
			// Prints out rd and rn
			_, err = fmt.Fprintf(f, "R%d, R%d, ", list[i].rd, list[i].rn)
			// Prints out shamt or rm depending on opcode
			if list[i].op == "LSL" || list[i].op == "ASR" || list[i].op == "LSR" {
				_, err = fmt.Fprintf(f, "#%d\n", list[i].shamt)
			} else {
				_, err = fmt.Fprintf(f, "R%d\n", list[i].rm)
			}

		case "D":
			_, err := fmt.Fprintf(f, "%s %s %s %s %s\t", list[i].rawInstruction[0:11],
				list[i].rawInstruction[11:20], list[i].rawInstruction[20:22],
				list[i].rawInstruction[22:27], list[i].rawInstruction[27:32])

			_, err = fmt.Fprintf(f, "%d\t%s ", list[i].programCnt, list[i].op)

			_, err = fmt.Fprintf(f, "R%d, [R%d,#%d]\n", list[i].rt, list[i].rn, list[i].address)
			if err != nil {
				log.Fatal(err)
			}

		case "I":
			_, err := fmt.Fprintf(f, "%s %s %s %s\t", list[i].rawInstruction[0:10],
				list[i].rawInstruction[10:22], list[i].rawInstruction[22:27],
				list[i].rawInstruction[27:32])

			_, err = fmt.Fprintf(f, "\t%d\t%s ", list[i].programCnt, list[i].op)

			_, err = fmt.Fprintf(f, "R%d, R%d, #%d\n", list[i].rd, list[i].rn, list[i].immediate)
			if err != nil {
				log.Fatal(err)
			}
		case "B":
			_, err := fmt.Fprintf(f, "%s %s\t", list[i].rawInstruction[0:6],
				list[i].rawInstruction[6:32])

			_, err = fmt.Fprintf(f, "\t%d\t%s\t", list[i].programCnt, list[i].op)

			_, err = fmt.Fprintf(f, "#%d\n", list[i].offset)
			if err != nil {
				log.Fatal(err)
			}
		case "CB":
			_, err := fmt.Fprintf(f, "%s %s %s\t", list[i].rawInstruction[0:8],
				list[i].rawInstruction[8:27], list[i].rawInstruction[27:32])

			_, err = fmt.Fprintf(f, "\t%d\t%s ", list[i].programCnt, list[i].op)

			_, err = fmt.Fprintf(f, "R%d #%d\n", list[i].conditional, list[i].offset)
			if err != nil {
				log.Fatal(err)
			}
		case "IM":
			_, err := fmt.Fprintf(f, "%s %s %s %s\t", list[i].rawInstruction[0:9],
				list[i].rawInstruction[9:11], list[i].rawInstruction[11:27],
				list[i].rawInstruction[27:32])

			_, err = fmt.Fprintf(f, "\t%d\t%s ", list[i].programCnt, list[i].op)

			_, err = fmt.Fprintf(f, "R%d, %d, LSL %d\n", list[i].rd, list[i].field, list[i].shiftCode)
			if err != nil {
				log.Fatal(err)
			}
		case "BREAK":
			_, err = fmt.Fprintf(f, "%s %s %s %s %s %s\t%d BREAK\n", list[i].rawInstruction[:8], list[i].rawInstruction[8:11], list[i].rawInstruction[11:16], list[i].rawInstruction[16:21],
				list[i].rawInstruction[21:26], list[i].rawInstruction[26:], list[i].programCnt)
			if err != nil {
				log.Fatal(err)
			}
		case "NOP":
			_, err = fmt.Fprintf(f, "%s \t\t%d %s\n", list[i].rawInstruction, list[i].programCnt, list[i].op)
			if err != nil {
				log.Fatal(err)
			}
		case "NUM":
			_, err = fmt.Fprintf(f, "%s \t\t%d %d\n", list[i].rawInstruction, list[i].programCnt, list[i].bitValue)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func TwoComplement(i uint64, bitLength uint) int64 {
	var n int64
	var v int64
	n = int64(i)
	v = (1 << bitLength) - 1
	if (i >> (bitLength - 1)) != 0 {
		n = ((n ^ v) + 1) * -1
	}
	return n
}

/*
***************************************
// Format type of Instruction
***************************************
*/
func RTypeFormat(list *Instruction) {

	// bits 12 - 16
	list.rm = uint8((list.linevalue & 2031616) >> 16)
	// bits 17 - 22
	list.shamt = uint8((list.linevalue & 64512) >> 10)
	// bits 23 - 27
	list.rn = uint8((list.linevalue & 992) >> 5)
	// bits 28 - 32
	list.rd = uint8(list.linevalue & 31)
}

func DTypeFormat(list *Instruction) {
	// bits 12 - 20
	list.address = uint16((list.linevalue & 2093056) >> 12)
	// bits 21 - 22
	list.op2 = uint8((list.linevalue & 3072) >> 10)
	// bits 23 - 27
	list.rn = uint8((list.linevalue & 992) >> 5)
	// bits 28 - 32
	list.rt = uint8(list.linevalue & 31)
}

func ITypeFormat(list *Instruction) {
	// bits 11 - 22
	list.immediate = int32(TwoComplement((list.linevalue&4193280)>>10, 12))
	// bits 23 - 27
	list.rn = uint8((list.linevalue & 992) >> 5)
	// bits 28 -32
	list.rd = uint8(list.linevalue & 31)
}

func BTypeFormat(list *Instruction) {
	// bits 7 - 32
	list.offset = int32(TwoComplement((list.linevalue & 67108863), 26))
}

func CBTypeFormat(list *Instruction) {
	// bits 9 - 27
	list.offset = int32(TwoComplement((list.linevalue&16777184)>>5, 19))
	// bits 28 - 32
	list.conditional = uint8(list.linevalue & 31)
}

func IMTypeFormat(list *Instruction) {
	// bits 10-11
	list.shiftCode = uint8((list.linevalue & 6291456) >> 21)
	// bits 12 - 27
	list.field = uint32((list.linevalue & 2097120) >> 5)
	// bits 28 - 32
	list.rd = uint8((list.linevalue & 31))
}

func Break(list *Instruction) {
	// bits 11 - 32
	list.brk = uint32(list.linevalue & 2031591)
	Breaknow = true
}
