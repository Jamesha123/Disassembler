package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Instruction struct{
	typeofInstruction string
	rawInstruction string
	linevalue uint64
	programCnt int
	opcode uint64
	op string
	rd uint8
	rn uint8
	rm uint8
	shamt uint8
	im string
}

func processInput( list []Instruction){
	for i:=0; i<len(list); i++ {
		opcodeMasking(&list[i])
		convertToInt(&list[i])
		opcodeMatching(&list[i])

	}
}

// // mask the opcode (11 bits)
// func opcodeMasking(ins *Instruction){

// }

// converts opcode to int
func convertToInt(instruct *Instruction){
	i, err := strconv.ParseUint(instruct.rawInstruction, 2, 64)
	if err == nil{
		instruct.linevalue = i;
	} else{
		fmt.println(err)
	}
}

// match the converted opcode to the arm instruction
func opcodeMatching(ins *Instruction){
	if ins.opcode >= 160 && ins.opcode <= 191 {
		ins.op = "B"
	} else if ins.opcode == 1104 {
		ins.op = "AND"
		ins.instructionType = "R"
	} else if ins.opcode == 1112 {
		ins.op = "ADD"
		ins.instructionType = "R"
	} else if ins.opcode >= 1160 && ins.opcode <= 1161 {
		ins.op = "ADDI"
		ins.instructionType = "I"
	} else if ins.opcode == 1360 {
		ins.op = "ORR"
		ins.instructionType = "R"
	} else if ins.opcode >= 1440 && ins.opcode <= 1447 {
		ins.op = "CBZ"
		ins.instructionType = "CB"
	} else if ins.opcode >= 1448 && ins.opcode <= 1455 {
		ins.op = "CBNZ"
		ins.instructionType = "CB"
	} else if ins.opcode == 1642 {
		ins.op = "SUB"
		ins.instructionType = "R"
	} else if ins.opcode >= 1672 && ins.opcode <= 1673 {
		ins.op = "SUBI"
		ins.instructionType = "I"
	} else if ins.opcode >= 1684 && ins.opcode <= 1687 {
		ins.op = "MOVZ"
		ins.instructionType = "IM"
	} else if ins.opcode >= 1940 && ins.opcode <= 1943 {
		ins.op = "MOVK"
		ins.instructionType = "IM"
	} else if ins.opcode == 1690 {
		ins.op = "LSR"
		ins.instructionType = "R"
	} else if ins.opcode == 1691 {
		ins.op = "LSL"
		ins.instructionType = "R"
	} else if ins.opcode == 1984 {
		ins.op = "STUR"
		ins.instructionType = "D"
	} else if ins.opcode == 1986 {
		ins.op = "LDUR"
		ins.instructionType = "D"
	} else if ins.opcode == 1692 {
		ins.op = "ASR"
		ins.instructionType = "R"
	} else if ins.opcode == 0 {
		ins.op = "NOP"
	} else if ins.opcode == 1872 {
		ins.op = "EOR"
		ins.instructionType = "R"
	} else if ins.opcode == 2038 {
		ins.op = "Break"
	} else {
		fmt.Println("Invalid opcode")
	}
}

func main() {

	fmt.Println("Opening a file")
	var file, err := os.Open("addtest1_bins.txt")
	if err != nil {
		log.Fatalf("Could not open file: ", err)
	}

	var Breaknow bool = false
	InputParsed := []Instruction{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan(){
		ns:= Instruction{rawInstruction: scanner.Text()}
		InputParsed = append(InputParsed, ns)
	}

	if err := scanner.Err(); err != nil{
		fmt.Println(err)
	}

	defer file.Close()

}

