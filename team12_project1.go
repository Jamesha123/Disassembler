package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Instruction struct {
	typeofInstruction string
	rawInstruction    string //binary input converted to string
	linevalue         uint64 //
	programCnt        int
	opcode            uint64 //needs to be 64 bit to apply mask and shift
	op                string //whether its B, I, BREAK, etc
	op2               uint8  //2 bits in D format
	rd                uint8  //5 bit, 0-31 registers
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
	bitValue          int64
}

type Snapshot struct {
	cycle int
	regis [32]int
	PC    int
}

/*
*************************************
// functions
*************************************
*/
func processInput(list []Instruction) {

	for i := 0; i < len(list); i++ {
		if !Breaknow {
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
				Break(&list[i])
				PC = list[i].programCnt + 4
			}
		} else {
			list[i].typeofInstruction = "NUM"
			var val uint64
			val, _ = strconv.ParseUint(list[i].rawInstruction, 2, 32)
			list[i].bitValue = TwoComplement(val, 32)
			counter := list[i].programCnt
			value := int(list[i].bitValue)
			temp := []int{counter, value}
			Data = append(Data, temp)
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
		ins.typeofInstruction = "B"
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
	} else if ins.opcode == 0 && ins.linevalue == 0 {
		ins.op = "NOP"
		ins.typeofInstruction = "NOP"
	} else if ins.opcode == 1872 {
		ins.op = "EOR"
		ins.typeofInstruction = "R"
	} else if ins.opcode == 2038 {
		ins.op = "Break"
		ins.typeofInstruction = "BREAK"
	} else {
		//fmt.Println("Error")
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
			_, err := fmt.Fprintf(f, "%s %s %s %s %s\t", list[i].rawInstruction[0:11],
				list[i].rawInstruction[11:20], list[i].rawInstruction[20:22],
				list[i].rawInstruction[22:27], list[i].rawInstruction[27:32])

			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].programCnt, list[i].op)

			_, err = fmt.Fprintf(f, "R%d, [R%d,#%d]\n", list[i].rt, list[i].rn, list[i].address)
			if err != nil {
				log.Fatal(err)
			}

		case "I":
			_, err := fmt.Fprintf(f, "%s %s %s %s\t", list[i].rawInstruction[0:10],
				list[i].rawInstruction[10:22], list[i].rawInstruction[22:27],
				list[i].rawInstruction[27:32])

			_, err = fmt.Fprintf(f, "\t%d\t%s\t", list[i].programCnt, list[i].op)

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

			_, err = fmt.Fprintf(f, "\t%d\t%s\t", list[i].programCnt, list[i].op)

			_, err = fmt.Fprintf(f, "R%d #%d\n", list[i].conditional, list[i].offset)
			if err != nil {
				log.Fatal(err)
			}
		case "IM":
			_, err := fmt.Fprintf(f, "%s %s %s %s\t", list[i].rawInstruction[0:9],
				list[i].rawInstruction[9:11], list[i].rawInstruction[11:27],
				list[i].rawInstruction[27:32])

			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].programCnt, list[i].op)

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
	//var n int64
	// for i, v := range bitLength {
	// 	shift := uint((len(bitLength) - i - 1) * 8)
	// 	if i == 0 && v&0x80 != 0 {
	// 		n -= 0x80 << shift
	// 		v &= 0x7f
	// 	}
	// 	n += int64(v) << shift
	// }
	//return n

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
func RTypeFormat(ins *Instruction) {

	// bits 12 - 16
	ins.rm = uint8((ins.linevalue & 2031616) >> 16)
	// bits 17 - 22
	ins.shamt = uint8((ins.linevalue & 64512) >> 10)
	// bits 23 - 27
	ins.rn = uint8((ins.linevalue & 992) >> 5)
	// bits 28 - 32
	ins.rd = uint8(ins.linevalue & 31)
}

func DTypeFormat(ins *Instruction) {
	// bits 12 - 20
	ins.address = uint16((ins.linevalue & 2093056) >> 12)
	// bits 21 - 22
	ins.op2 = uint8((ins.linevalue & 3072) >> 10)
	// bits 23 - 27
	ins.rn = uint8((ins.linevalue & 992) >> 5)
	// bits 28 - 32
	ins.rt = uint8(ins.linevalue & 31)
}

func ITypeFormat(ins *Instruction) {
	// bits 11 - 22
	ins.immediate = int32(TwoComplement((ins.linevalue&4193280)>>10, 12))
	// bits 23 - 27
	ins.rn = uint8((ins.linevalue & 992) >> 5)
	// bits 28 -32
	ins.rd = uint8(ins.linevalue & 31)
}

func BTypeFormat(ins *Instruction) {
	// bits 7 - 32
	ins.offset = int32(TwoComplement((ins.linevalue & 67108863), 26))
}

func CBTypeFormat(ins *Instruction) {
	// bits 9 - 27
	ins.offset = int32(TwoComplement((ins.linevalue&16777184)>>5, 19))
	// bits 28 - 32
	ins.conditional = uint8(ins.linevalue & 31)
}

func IMTypeFormat(ins *Instruction) {
	// bits 10-11
	ins.shiftCode = uint8((ins.linevalue & 6291456) >> 21)
	// bits 12 - 27
	ins.field = uint32((ins.linevalue & 2097120) >> 5)
	// bits 28 - 32
	ins.rd = uint8((ins.linevalue & 31))
}

func Break(ins *Instruction) {
	// bits 11 - 32
	ins.brk = uint32(ins.linevalue & 2031591)
	Breaknow = true
}

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

func ProcessSnapshot(list []Instruction, array []Snapshot) {
	for i := 0; i < len(list); i++ {
		//if !Breaknow {
		if i != 0 {
			array[i].regis = array[i-1].regis //copies previous cycle's registers to keep track of register status
		}

		// switch list[i].typeofInstruction {
		// case "R":
		// 	RTypeArithmetic(&list[i], &array[i])
		// case "D":
		// 	DTypeArithmetic(&list[i], &array[i])
		// case "I":
		// 	ITypeArithmetic(&list[i], &array[i])
		// //case "B":
		// //	BTypeArithmetic(&list[i], &array[i])
		// case "CB":
		// 	CBTypeArithmetic(&list[i], &array[i])
		// case "IM":
		// 	IMTypeArithmetic(&list[i], &array[i])
		// }
	}
}

func RTypeArithmetic(list *Instruction, array *Snapshot) {
	switch list.op {
	case "ADD":
		array.regis[list.rd] = array.regis[list.rm] + array.regis[list.rn]
	case "SUB":
		array.regis[list.rd] = array.regis[list.rm] - array.regis[list.rn]
	case "AND":
		array.regis[list.rd] = array.regis[list.rm] & array.regis[list.rn]
	case "ORR":
		array.regis[list.rd] = array.regis[list.rm] | array.regis[list.rn]
	case "EOR":
		array.regis[list.rd] = array.regis[list.rm] ^ array.regis[list.rn]
	case "LSL":
		array.regis[list.rd] = array.regis[list.rn] << list.shamt
	case "ASR":
		array.regis[list.rd] = array.regis[list.rn] >> list.shamt
	case "LSR":
		array.regis[list.rd] = array.regis[list.rn] >> uint(list.shamt)
	}

}

func DTypeArithmetic(list *Instruction, array *Snapshot) {
	switch list.op {
	case "LDUR":
		//array.regis[list.rt] = array.regis[list.rn] + array.regis[list.address]

	case "STUR":
		var inSlice bool
		counter := array.regis[list.rn] + int(list.address) // + address = 100
		//inSlice = counterInSlice(counter, SturData)
		if !inSlice {
			value := array.regis[list.rt]
			temp := []int{counter, value}
			SturData = append(SturData, temp)
		}
	}
}

func ITypeArithmetic(list *Instruction, array *Snapshot) {
	switch list.op {
	case "ADDI":
		array.regis[list.rd] = array.regis[list.rn] + int(list.immediate)
	case "SUBI":
		array.regis[list.rd] = array.regis[list.rn] - int(list.immediate)
	}
}

func BTypeArithmetic(list *Instruction, array *Snapshot) {
	array.regis[list.programCnt] += int(list.offset * 4)
}

func IMTypeArithmetic(list *Instruction, array *Snapshot) {
	switch list.op {
	case "MOVZ":
		array.regis[list.rd] = int(list.field) << int(list.shiftCode*16)
	case "MOVK":
		array.regis[list.rd] = int(uint(array.regis[list.rd])) ^ int(list.address)<<(int(list.shamt)*16)
	}
}

func CBTypeArithmetic(list *Instruction, array *Snapshot) {
	switch list.op {
	case "CBZ":
		if array.regis[list.conditional] == 0 {
			list.programCnt += int(list.offset * 4)
		}

	case "CBNZ":
		if array.regis[list.conditional] != 0 {
			list.programCnt += int(list.offset * 4)
		}
	}
}

func writeSimulator(filePath string, list []Instruction, array []Snapshot) {
	f, err := os.Create(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	// Breaknow = false
	// line := "===================="
	// cycleLabel := "cycle:"
	// registersLabel := "registers:"
	// registerNumberlabel := []string{"r00", "r08", "r16", "r24"}
	dataLabel := "data:"

	for i := 0; i < len(list); i++ {
		// fmt.Println(list[i].programCnt)
		// _, err := fmt.Fprintf(f, "%s\n", line)
		// switch list[i].typeofInstruction {
		// case "R":
		// 	_, err = fmt.Fprintf(f, "%s%d\t%d\t%s R%d, R%d, ", cycleLabel, array[i].cycle, array[i].PC, list[i].op, list[i].rd, list[i].rn)
		// 	if list[i].op == "LSL" || list[i].op == "ASR" || list[i].op == "LSR" {
		// 		_, err = fmt.Fprintf(f, "#%d\n", list[i].shamt)
		// 	} else {
		// 		_, err = fmt.Fprintf(f, "R%d\n", list[i].rm)
		// 	}
		// case "D":
		// 	_, err = fmt.Fprintf(f, "%s%d\t%d\t%s R%d, [R%d,#%d]\n", cycleLabel, array[i].cycle, list[i].programCnt, list[i].op, list[i].rt, list[i].rn, list[i].address)
		// case "I":
		// 	_, err = fmt.Fprintf(f, "%s%d\t%d\t%s R%d, R%d, #%d\n", cycleLabel, array[i].cycle, list[i].programCnt, list[i].op, list[i].rd, list[i].rn, list[i].immediate)
		// case "B":
		// 	_, err = fmt.Fprintf(f, "%s%d\t%d\t%s #%d\n", cycleLabel, array[i].cycle, list[i].programCnt, list[i].op, list[i].offset)
		// case "CB":
		// 	_, err = fmt.Fprintf(f, "%s%d\t%d\t%s R%d #%d\n", cycleLabel, array[i].cycle, list[i].programCnt, list[i].op, list[i].conditional, list[i].offset)
		// case "IM":
		// 	_, err = fmt.Fprintf(f, "%s%d\t%d\t%s R%d, %d, LSL %d\n", cycleLabel, array[i].cycle, list[i].programCnt, list[i].op, list[i].rd, list[i].field, list[i].shiftCode)
		// case "NOP":
		// 	_, err = fmt.Fprintf(f, "%s%d\t%d\t%s\n", cycleLabel, array[i].cycle, list[i].programCnt, list[i].op)
		// case "BREAK":
		// 	_, err = fmt.Fprintf(f, "%s%d\t%d\t%s\n", cycleLabel, array[i].cycle, list[i].programCnt, list[i].op)
		// 	Breaknow = true
		// 	//case "NUM":
		// 	//	_, err = fmt.Fprintf(f, "%s%d\t%d\t%d\n", cycleLabel, array[i].cycle, list[i].programCnt, list[i].bitValue)
		// }
		// _, err = fmt.Fprintf(f, "\n%s\n", registersLabel)

		// //_, err = fmt.Fprintf(f, "%s:\t", registerNumberlabel[0])

		// _, err = fmt.Fprintf(f, "%s:\t%5d\t%5d\t%5d\t%5d\t%5d\t%5d\t%5d\t%5d\n", registerNumberlabel[0], array[i].regis[0], array[i].regis[1], array[i].regis[2], array[i].regis[3], array[i].regis[4], array[i].regis[5], array[i].regis[6], array[i].regis[7])
		// _, err = fmt.Fprintf(f, "%s:\t%5d\t%5d\t%5d\t%5d\t%5d\t%5d\t%5d\t%5d\n", registerNumberlabel[1], array[i].regis[8], array[i].regis[9], array[i].regis[10], array[i].regis[11], array[i].regis[12], array[i].regis[13], array[i].regis[14], array[i].regis[15])
		// _, err = fmt.Fprintf(f, "%s:\t%5d\t%5d\t%5d\t%5d\t%5d\t%5d\t%5d\t%5d\n", registerNumberlabel[2], array[i].regis[16], array[i].regis[17], array[i].regis[18], array[i].regis[19], array[i].regis[20], array[i].regis[21], array[i].regis[22], array[i].regis[23])
		// _, err = fmt.Fprintf(f, "%s:\t%5d\t%5d\t%5d\t%5d\t%5d\t%5d\t%5d\t%5d\n", registerNumberlabel[3], array[i].regis[24], array[i].regis[25], array[i].regis[26], array[i].regis[27], array[i].regis[28], array[i].regis[29], array[i].regis[30], array[i].regis[31])

		_, err = fmt.Fprintf(f, "\n%s\n", dataLabel)

		mod := PC % 32
		var counter int
		var value int
		loopCounter := 0

		for i := range Data {
			loopCounter++
			counter = Data[i][0]
			value = Data[i][1]
			if loopCounter%9 != 0 {
				if counter%32 == mod {
					_, err = fmt.Fprintf(f, "%d:", counter)
					_, err = fmt.Fprintf(f, "%d\t", value)
				} else {
					_, err = fmt.Fprintf(f, "%d\t", value)
				}
			} else {
				_, err = fmt.Fprintf(f, "\n") //\n for last printed data
				mod := Data[i][0] % 32
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
		for i := range SturData {
			loopCounter2++
			counter2 = SturData[i][0]
			value2 = SturData[i][1]
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
				mod := SturData[i][0] % 32
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

		if list[i].typeofInstruction == "BREAK" {
			break
		}

		PC++
	}
}

// func counterInSlice(c int, slice [][]int) bool {
// 	for i := range SturData {
// 		if c == SturData[i][0] {
// 			return true
// 		}
// 	}
// 	return false
// }

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/*
********************************************
//  main
********************************************
*/

var PC = 96
var Data [][]int
var SturData [][]int
var InputParsed []Instruction
var SnapshotArray []Snapshot
var Breaknow bool = false
var InputFileName *string
var OutputFileName *string
var OutputFileName2 *string

func main() {

	//team := "team12_out"
	dis := "_dis.txt"
	sim := "_sim.txt"
	//InputFileName := flag.String("i", "dtest2_bin.txt", "Gets the input file name")
	InputFileName := flag.String("i", "dtest2_bin.txt", "Gets the input file name")
	OutputFileName := flag.String("o", "team12_out", "Gets the output file name")
	OutputFileName2 := flag.String("k", "team12_out", "Gets the output file name")

	flag.Parse()

	readInstruction(*InputFileName)
	processInput(InputParsed)
	initializeRegisters(InputParsed)
	ProcessSnapshot(InputParsed, SnapshotArray)
	writeInstruction(*OutputFileName+dis, InputParsed)

	writeSimulator(*OutputFileName2+sim, InputParsed, SnapshotArray)

	fmt.Println("end project 1")
}
