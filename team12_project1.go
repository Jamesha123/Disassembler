package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Instruction struct {
	typeofInstruction string
	rawInstruction    string
	linevalue         uint64
	programCnt        int
	opcode            uint64
	op                string
	rd                uint8
	rn                uint8
	rm                uint8
	shamt             uint8
	im                string
}

/*
*************************************
// functions
*************************************
*/
func processInput(list []Instruction) {
	for i := 0; i < len(list); i++ {

		convertToInt(&list[i])
		opcodeMasking(&list[i])
		opcodeMatching(&list[i])
		switch list[i].typeofInstruction {
		case "R":
			RTypeFormat(&list[i])
		case "D":

		case "I":

		case "B":

		case "CB":

		case "IM":
		}
	}
}

// mask the opcode (11 bits)
func opcodeMasking(ins *Instruction) {
	// mask the bits by 0xFFE00000 = 4292870144 to get first 11 bits
	ins.opcode = (ins.linevalue & 4292870144) >> 21
}

// converts opcode to int
func convertToInt(ins *Instruction) {
	i, err := strconv.ParseUint(ins.rawInstruction, 2, 64)
	//fmt.Println(i)
	if err == nil {
		ins.linevalue = i
	} else {
		fmt.Println(err)
	}
}

// match the converted opcode to the arm instruction
func opcodeMatching(ins *Instruction) {
	if ins.opcode >= 160 && ins.opcode <= 191 {
		ins.op = "B"
	} else if ins.opcode == 1104 {
		ins.op = "AND"
		ins.typeofInstruction = "R"
	} else if ins.opcode == 1112 {
		ins.op = "ADD"
		ins.typeofInstruction = "R"
	} else if ins.opcode >= 1160 && ins.opcode <= 1161 {
		ins.op = "ADDI"
		ins.typeofInstruction = "I"
	} else if ins.opcode == 1360 {
		ins.op = "ORR"
		ins.typeofInstruction = "R"
	} else if ins.opcode >= 1440 && ins.opcode <= 1447 {
		ins.op = "CBZ"
		ins.typeofInstruction = "CB"
	} else if ins.opcode >= 1448 && ins.opcode <= 1455 {
		ins.op = "CBNZ"
		ins.typeofInstruction = "CB"
	} else if ins.opcode == 1624 {
		ins.op = "SUB"
		ins.typeofInstruction = "R"
	} else if ins.opcode >= 1672 && ins.opcode <= 1673 {
		ins.op = "SUBI"
		ins.typeofInstruction = "I"
	} else if ins.opcode >= 1684 && ins.opcode <= 1687 {
		ins.op = "MOVZ"
		ins.typeofInstruction = "IM"
	} else if ins.opcode >= 1940 && ins.opcode <= 1943 {
		ins.op = "MOVK"
		ins.typeofInstruction = "IM"
	} else if ins.opcode == 1690 {
		ins.op = "LSR"
		ins.typeofInstruction = "R"
	} else if ins.opcode == 1691 {
		ins.op = "LSL"
		ins.typeofInstruction = "R"
	} else if ins.opcode == 1984 {
		ins.op = "STUR"
		ins.typeofInstruction = "D"
	} else if ins.opcode == 1986 {
		ins.op = "LDUR"
		ins.typeofInstruction = "D"
	} else if ins.opcode == 1692 {
		ins.op = "ASR"
		ins.typeofInstruction = "R"
	} else if ins.opcode == 0 {
		ins.op = "NOP"
	} else if ins.opcode == 1872 {
		ins.op = "EOR"
		ins.typeofInstruction = "R"
	} else if ins.opcode == 2038 {
		ins.op = "Break"
	} else {
		fmt.Println("Invalid opcode")
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
	//// InputParsed := []Instruction{}
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
			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].programCnt, list[i].op)
			// Prints out rd and rn
			_, err = fmt.Fprintf(f, "R%d, R%d, ", list[i].rd, list[i].rn)
			// Prints out shamt or rm depending on opcode
			if list[i].op == "LSL" || list[i].op == "ASR" || list[i].op == "LSR" {
				_, err = fmt.Fprintf(f, "#%d\n", list[i].shamt)
			} else {
				_, err = fmt.Fprintf(f, "R%d\n", list[i].rm)
			}

		case "D":
			// TO DO
		case "I":
			// TO DO
		case "B":
			// TO DO
		case "CB":
			// TO DO
		case "IM":
			// TO DO
		}
	}
}

/*
***************************************
// Format type of Instruction
***************************************
*/
func RTypeFormat(ins *Instruction) {

	// bits 12-16
	ins.rm = uint8((ins.linevalue & 2031616) >> 16)
	// bits 17-22
	ins.shamt = uint8((ins.linevalue & 64512) >> 10)
	// bits 23-27
	ins.rn = uint8((ins.linevalue & 992) >> 5)
	// bits 28-32
	ins.rd = uint8(ins.linevalue & 31)

}

func DTypeFormat(ins *Instruction) {
	// TO DO
}

func ITypeFormat(ins *Instruction) {
	// TO DO
}

func BTypeFormat(ins *Instruction) {
	// TO DO
}

func CBTypeFormat(ins *Instruction) {
	// TO DO
}

func IMTypeFormat(ins *Instruction) {
	// TO DO
}

/*
********************************************
//  main
********************************************
*/

var InputParsed []Instruction

func main() {

	os.Args[0] = "addtest1_bin.txt"

	readInstruction(os.Args[0])
	processInput(InputParsed)
	writeInstruction("team12_addtest1_out.dis.txt", InputParsed)

}
