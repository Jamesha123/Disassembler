package main

import "fmt"

var ALUVal int

func ALUOp(operand1 int, operand2 int, op string) {
	switch op {
	case "AND":
		ALUVal = operand1 & operand2
	case "ADD":
		ALUVal = operand1 + operand2
	case "ORR":
		ALUVal = operand1 | operand2
	case "SUB":
		ALUVal = operand1 - operand2
	case "LSR":
		ALUVal = operand1 >> operand2
	case "LSL":
		ALUVal = operand1 << operand2
	case "ASR":
		ALUVal = operand1 >> operand2
	case "EOR":
		ALUVal = operand1 ^ operand2
	}
}

func ALUOpImm(immediate int, operand int, op string) {
	if op == "ADDI" {
		ALUVal = operand + immediate
	} else if op == "SUBI" {
		ALUVal = operand + immediate
	} else {
		fmt.Println("error")
	}
}

func ALUCall(list Instruction) int {
	if list.typeofInstruction == "I" {
		ALUOpImm(Register[list.rn], int(list.immediate), list.op)
	} else {
		ALUOp(Register[list.rn], Register[list.rm], list.op)
	}
	return ALUVal
}
