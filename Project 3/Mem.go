package main

func MEM(ins Instruction) int {
	if ins.op == "LDUR" {

		return LDURMem(Register[ins.rn] + int(ins.address)*4)

	} else if ins.op == "STUR" {

		STURMem(Register[ins.rn]+int(ins.address)*4, Register[ins.rt])

	}
	return -1
}
