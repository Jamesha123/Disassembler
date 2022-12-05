package main

func WriteBack(list Instruction, value int) {
	if list.typeofInstruction == "D" {
		Register[list.rt] = value
	} else {
		Register[list.rd] = value
	}
}
