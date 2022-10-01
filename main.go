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

func processInput(list []Instruction) {
	for i := 0; i < len(list); i++ {
		//opcodeMasking(&list[i])
		convertToInt(&list[i])
		opcodeMatching(&list[i])

	}
}

// // mask the opcode (11 bits)
// func opcodeMasking(ins *Instruction){

// }

// converts opcode to int
func convertToInt(ins *Instruction) {
	i, err := strconv.ParseUint(ins.rawInstruction, 2, 64)
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
	} else if ins.opcode == 1642 {
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

func main() {

	fmt.Println("Opening a file")
	file, err := os.Open("addtest1_bin.txt")
	if err != nil {
		log.Fatalf("Could not open file: ", err)
	}

	//var Breaknow bool = false
	InputParsed := []Instruction{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ns := Instruction{rawInstruction: scanner.Text()}
		InputParsed = append(InputParsed, ns)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	defer file.Close()

}
