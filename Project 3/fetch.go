package main

var pc int = 0

func Fetch() bool {
	list := InputParsed[pc]
	cacheHit, _ := CheckCacheHit(int(list.bitValue))
	branch := false
	if cacheHit {
		if list.op == "B" {
			pc = pc + int(list.offset)
			branch = true
		} else if list.op == "CBZ" {
			if Register[list.conditional] == 0 {
				pc = pc + int(list.offset)
				branch = true
			}
		} else if list.op == "CBNZ" {
			if Register[list.conditional] != 0 {
				pc = pc + int(list.offset)
				branch = true
			}
		} else if list.op == "BREAk" {
			Breaknow = true
		} else if list.op == "NOP" {

		} else {
			PreIssueBuff <- pc
			pc++
		}
	} else {

	}
	return branch
}
