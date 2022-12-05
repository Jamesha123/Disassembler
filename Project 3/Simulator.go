package main

import (
	"fmt"
	"log"
	"os"
)

type Snapshot struct {
	cycle int
	regis [32]int
	PC    int
}

// variables
var Register [32]int
var BreakPoint int
var PCIndex = 0
var Data [][]int
var SturData [][]int
var InputParsed []Instruction
var SnapshotArray []Snapshot
var Breaknow bool = false
var InputFileName *string
var OutputFileName *string
var OutputFileName2 *string
var PC = 96
var MemoryIndex = make(map[int]int)
var PreIssueBuff = make(chan int, 4)
var PreMemBuff = make(chan int, 2)
var PreALUBuff = make(chan int, 2)
var postMemBuff = make(chan [2]int, 1)
var postALUBuff = make(chan [2]int, 1)
var cycleNum = 0

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func initializeRegisters(list []Instruction) {
	round := 1
	//initReg := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	initReg := [32]int{5, 8, 200, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < len(list); i++ {
		temp := Snapshot{
			cycle: round,
		}
		SnapshotArray = append(SnapshotArray, temp)
		SnapshotArray[i].regis = initReg
		round++
	}
}

func ProcessSnapshot(list []Instruction) {
	for i := 0; i < len(list); i++ {
		//if !Breaknow {
		if i != 0 {
			Register[i] = Register[i-1] //copies previous cycle's Register to keep track of register status
		}
	}
}

func ExecuteInstruction(list Instruction) {
	switch list.op {
	case "ADD":
		Register[list.rd] = Register[list.rn] + Register[list.rm]
	case "SUB":
		Register[list.rd] = Register[list.rn] - Register[list.rm]
	case "AND":
		Register[list.rd] = Register[list.rn] & Register[list.rm]
	case "ORR":
		Register[list.rd] = Register[list.rn] | Register[list.rm]
	case "EOR":
		Register[list.rd] = Register[list.rn] ^ Register[list.shamt]
	case "LSL":
		Register[list.rd] = Register[list.rn] << Register[list.shamt]
	case "ASR":
		Register[list.rd] = Register[list.rn] >> Register[list.shamt]
	case "LSR":
		Register[list.rd] = Register[list.rn] >> Register[list.shamt]
	case "LDUR":
		counter := Register[list.rn] + int(list.address)
		inSlice := counterInSlice(counter, SturData)
		if !inSlice {
			Register[list.rt] = int(MemoryIndex[Register[list.rn]+int(list.address)*4])
			//fmt.Println("register = ", Register[list.rt])
		} else {
			fmt.Println("Out of Memory")
		}
	case "STUR":
		var inSlice bool
		counter := Register[list.rn] + int(list.address) // + address = 100
		inSlice = counterInSlice(counter, SturData)
		if !inSlice {
			MemoryIndex[Register[list.rn]+int(list.address)*4] = Register[list.rt]
			// value := Register[list.rt]
			// temp := []int{counter, value}
			// SturData = append(SturData, temp)
			// fmt.Println("val =", value, "temp =", temp, counter = ", counter, "SturData = ", SturData, "rn =", list.rm, "address = ", list.address)
		} else {
			fmt.Println("Out of Memory")
		}
	case "ADDI":
		Register[list.rd] = Register[list.rn] + int(list.immediate)
	case "SUBI":
		Register[list.rd] = Register[list.rn] - int(list.immediate)
	case "B":
		PCIndex += int(list.offset)
	case "MOVZ":
		Register[list.rd] = int(list.field) << 16 * int(list.shiftCode)
	case "MOVK":
		Register[list.rd] = int(uint16(Register[list.rd])) ^ int(list.address)<<(list.shamt*16)
	case "CBZ":
		if Register[list.conditional] == 0 {
			PCIndex += int(list.offset)
		}
	case "CBNZ":
		if Register[list.conditional] != 0 {
			PCIndex += int(list.offset)
		}
	case "NOP":
		break
	}
}

func SimulateCycle() {
	//Write back with both postALUBuff and postMemBuff
	if len(postALUBuff) != 0 {
		buff := <-postALUBuff
		WriteBack(InputParsed[buff[1]], buff[0])
	}
	if len(postMemBuff) != 0 {
		buff := <-postMemBuff
		WriteBack(InputParsed[buff[1]], buff[0])
	}

	if len(PreALUBuff) != 0 {
		listIndex := <-PreALUBuff
		var ALUOut = [2]int{listIndex, ALUCall(InputParsed[listIndex])}
		postALUBuff <- ALUOut
	}

	if len(PreMemBuff) != 0 {
		listIndex := <-PreMemBuff
		//check for cache hit
		cacheHit, _ := CheckCacheHit(Register[InputParsed[listIndex].rn] + int(InputParsed[listIndex].address)*4)
		if cacheHit {
			var MemOut = [2]int{listIndex, MEM(InputParsed[listIndex])}
			if MemOut[1] != -1 {
				postMemBuff <- MemOut
			}
		} else {
			// change MEM and cache
			MEM(InputParsed[listIndex])
			// put Index back in buffer
			if len(PreMemBuff) == 0 {
				PreMemBuff <- listIndex
			} else if len(PreMemBuff) == 1 {
				tempInt := <-PreMemBuff
				PreMemBuff <- listIndex
				PreMemBuff <- tempInt
			}
		}
	}

	if !Breaknow {
		for i := 0; i < 2; i++ {
			if Fetch() {
				break
			}
		}
	}

	if len(PreIssueBuff) != 0 {
		// Issue()
	}
}

func writeSimulator(filePath string, list []Instruction) {
	f, err := os.Create(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	var cycle = 1
	line := "===================="
	cycleLabel := "cycle:"
	registersLabel := "Register:"
	dataLabel := "data:"
	//fmt.Println(PCIndex)
	//fmt.Println(BreakPoint)

	for PCIndex < BreakPoint {
		_, err := fmt.Fprintf(f, "%s\n", line)
		switch list[PCIndex].typeofInstruction {
		case "R":
			_, err = fmt.Fprintf(f, "%s%d\t%d\t\t%s R%d, R%d, ", cycleLabel, cycle, list[PCIndex].programCnt, list[PCIndex].op, list[PCIndex].rd, list[PCIndex].rn)
			if list[PCIndex].op == "LSL" || list[PCIndex].op == "ASR" || list[PCIndex].op == "LSR" {
				_, err = fmt.Fprintf(f, "#%d\n", list[PCIndex].shamt)
			} else {
				_, err = fmt.Fprintf(f, "R%d\n", list[PCIndex].rm)
			}
		case "D":
			_, err = fmt.Fprintf(f, "%s%d\t%d\t\t%s R%d, [R%d,#%d]\n", cycleLabel, cycle, list[PCIndex].programCnt, list[PCIndex].op, list[PCIndex].rt, list[PCIndex].rn, list[PCIndex].address)
		case "I":
			_, err = fmt.Fprintf(f, "%s%d\t%d\t\t%s R%d, R%d, #%d\n", cycleLabel, cycle, list[PCIndex].programCnt, list[PCIndex].op, list[PCIndex].rd, list[PCIndex].rn, list[PCIndex].immediate)
		case "B":
			_, err = fmt.Fprintf(f, "%s%d\t%d\t\t%s #%d\n", cycleLabel, cycle, list[PCIndex].programCnt, list[PCIndex].op, list[PCIndex].offset)
		case "CB":
			_, err = fmt.Fprintf(f, "%s%d\t%d\t\t%s R%d #%d\n", cycleLabel, cycle, list[PCIndex].programCnt, list[PCIndex].op, list[PCIndex].conditional, list[PCIndex].offset)
		case "IM":
			_, err = fmt.Fprintf(f, "%s%d\t%d\t\t%s R%d, %d, LSL %d\n", cycleLabel, cycle, list[PCIndex].programCnt, list[PCIndex].op, list[PCIndex].rd, list[PCIndex].field, list[PCIndex].shiftCode)
		case "NOP":
			_, err = fmt.Fprintf(f, "%s%d\t%d\t\t%s\n", cycleLabel, cycle, list[PCIndex].programCnt, list[PCIndex].op)
		case "BREAK":
			_, err = fmt.Fprintf(f, "%s%d\t%d\t\t%s\n", cycleLabel, cycle, list[PCIndex].programCnt, list[PCIndex].op)
		}

		ExecuteInstruction(list[PCIndex])

		_, err = fmt.Fprintf(f, "\n%s\n", registersLabel)

		_, err = fmt.Fprintf(f, "r00:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", Register[0], Register[1], Register[2], Register[3], Register[4], Register[5], Register[6], Register[7])
		_, err = fmt.Fprintf(f, "r08:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", Register[8], Register[9], Register[10], Register[11], Register[12], Register[13], Register[14], Register[15])
		_, err = fmt.Fprintf(f, "r16:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", Register[16], Register[17], Register[18], Register[19], Register[20], Register[21], Register[22], Register[23])
		_, err = fmt.Fprintf(f, "r24:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n\n", Register[24], Register[25], Register[26], Register[27], Register[28], Register[29], Register[30], Register[31])

		_, err = fmt.Fprintf(f, "\n%s\n", dataLabel)

		mod := PC % 32
		var counter int
		var value int
		loopCounter := 0

		for PCIndex := range Data {
			loopCounter++
			counter = Data[PCIndex][0]
			value = Data[PCIndex][1]
			if loopCounter%9 != 0 {
				if counter%32 == mod {
					_, err = fmt.Fprintf(f, "%d:", counter)
					_, err = fmt.Fprintf(f, "%d\t", value)
				} else {
					_, err = fmt.Fprintf(f, "%d\t", value)
				}
			} else {
				_, err = fmt.Fprintf(f, "\n") //\n for last printed data
				mod := Data[PCIndex][0] % 32
				if counter%32 == mod {
					_, err = fmt.Fprintf(f, "%d:", counter)
					_, err = fmt.Fprintf(f, "%d\t", value)
				} else {
					_, err = fmt.Fprintf(f, "%d\t", value)
				}
			}
		}

		_, err = fmt.Fprintf(f, "\n")

		var mod2 int
		var counter2 int
		var value2 int
		loopCounter2 := 0
		for PCIndex := range SturData {
			loopCounter2++
			counter2 = SturData[PCIndex][0]
			value2 = SturData[PCIndex][1]
			mod2 = counter2 % 32
			if loopCounter2%9 != 0 {
				if counter2%32 == mod2 {
					_, err = fmt.Fprintf(f, "%d:", counter2)
					_, err = fmt.Fprintf(f, "%d\t", value)
				} else {
					_, err = fmt.Fprintf(f, "%d\t", value2)
				}
			} else {
				_, err = fmt.Fprintf(f, "\n") //\n for last printed data
				mod := SturData[PCIndex][0] % 32
				if counter2%32 == mod {
					_, err = fmt.Fprintf(f, "%d:", counter2)
					_, err = fmt.Fprintf(f, "%d\t", value2)
				} else {
					_, err = fmt.Fprintf(f, "%d\t", value2)
				}
			}
		}
		_, err = fmt.Fprintf(f, "\n") //\n for last printed data

		if err != nil {
			log.Fatal(err)
		}

		if list[PCIndex].typeofInstruction == "BREAK" {
			break
		}

		PCIndex++
		cycle++
	}

	for len(postALUBuff) != 0 || len(postMemBuff) != 0 || len(PreALUBuff) != 0 || len(PreMemBuff) != 0 {
		SimulateCycle()
		cycleNum++
	}
}

func counterInSlice(c int, slice [][]int) bool {
	for i := range SturData {
		if c == SturData[i][0] {
			return true
		}
	}
	return false
}
