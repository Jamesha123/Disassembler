package main

import (
	"flag"
	"fmt"
)

type Instruction struct {
	typeofInstruction string
	rawInstruction    string //binary input converted to string
	linevalue         uint64 //
	programCnt        int
	opcode            uint64 //needs to be 64 bit to apply mask and shift
	op                string //whether its B, I, BREAK, etc
	op2               uint8  //2 bits in D format
	rd                uint8  //5 bit, 0-31 Register
	rn                uint8
	rm                uint8
	rt                uint8
	shamt             uint8 //6 bits, 0-63
	im                string
	immediate         int32 //
	offset            int32 //
	conditional       uint8
	address           uint16
	shiftCode         uint8
	field             uint32
	brk               uint32
	bitValue          int
}

func main() {

	//InputFileName := flag.String("i", "dtest2_bin.txt", "Gets the input file name")
	OutputFileName := flag.String("o", "team12_out", "Gets the output file name")

	flag.Parse()

	readInstruction("dtest2_bin.txt")
	processInput(InputParsed)
	initializeRegisters(InputParsed)
	ProcessSnapshot(InputParsed)
	writeInstruction(*OutputFileName+"_dis.txt", InputParsed)
	writeSimulator(*OutputFileName+"_sim.txt", InputParsed)

	fmt.Println("end project")
}
