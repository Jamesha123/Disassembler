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
			BreakPoint = i
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
				Breaknow = true
				PC = list[i].programCnt + 4
				BreakPoint = i
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

		MemoryIndex[int(list[i].programCnt)] = list[i].bitValue
	}

}

// mask the opcode (11 bits)
func opcodeMasking(list *Instruction) {
	// mask the bits by 0xFFE00000 = 4292870144 to get first 11 bits
	list.opcode = (list.linevalue & 4292870144) >> 21
}

// converts opcode to int
func convertToInt(list *Instruction) {
	i, err := strconv.ParseUint(list.rawInstruction, 2, 64)
	//fmt.Println(i)
	if err == nil {
		list.linevalue = i
	} else {
		fmt.Println(err)
	}
}

// match the converted opcode to the arm instruction
func opcodeMatching(list *Instruction) {
	if list.opcode >= 160 && list.opcode <= 191 {
		list.op = "B"
		list.typeofInstruction = "B"
	} else if list.opcode == 1104 {
		list.op = "AND"
		list.typeofInstruction = "R"
	} else if list.opcode == 1112 {
		list.op = "ADD"
		list.typeofInstruction = "R"
	} else if list.opcode >= 1160 && list.opcode <= 1161 {
		list.op = "ADDI"
		list.typeofInstruction = "I"
	} else if list.opcode == 1360 {
		list.op = "ORR"
		list.typeofInstruction = "R"
	} else if list.opcode >= 1440 && list.opcode <= 1447 {
		list.op = "CBZ"
		list.typeofInstruction = "CB"
	} else if list.opcode >= 1448 && list.opcode <= 1455 {
		list.op = "CBNZ"
		list.typeofInstruction = "CB"
	} else if list.opcode == 1624 {
		list.op = "SUB"
		list.typeofInstruction = "R"
	} else if list.opcode >= 1672 && list.opcode <= 1673 {
		list.op = "SUBI"
		list.typeofInstruction = "I"
	} else if list.opcode >= 1684 && list.opcode <= 1687 {
		list.op = "MOVZ"
		list.typeofInstruction = "IM"
	} else if list.opcode >= 1940 && list.opcode <= 1943 {
		list.op = "MOVK"
		list.typeofInstruction = "IM"
	} else if list.opcode == 1690 {
		list.op = "LSR"
		list.typeofInstruction = "R"
	} else if list.opcode == 1691 {
		list.op = "LSL"
		list.typeofInstruction = "R"
	} else if list.opcode == 1984 {
		list.op = "STUR"
		list.typeofInstruction = "D"
	} else if list.opcode == 1986 {
		list.op = "LDUR"
		list.typeofInstruction = "D"
	} else if list.opcode == 1692 {
		list.op = "ASR"
		list.typeofInstruction = "R"
	} else if list.opcode == 0 && list.linevalue == 0 {
		list.op = "NOP"
		list.typeofInstruction = "NOP"
	} else if list.opcode == 1872 {
		list.op = "EOR"
		list.typeofInstruction = "R"
	} else if list.opcode == 2038 {
		list.op = "BREAK"
		list.typeofInstruction = "BREAK"
	} else {
		fmt.Println("Invalid Opcode")
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
	//InputParsed := []Instruction{}
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
			_, err = fmt.Fprintf(f, "%d\t%s ", list[i].programCnt, list[i].op)
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

			_, err = fmt.Fprintf(f, "%d\t%s ", list[i].programCnt, list[i].op)

			_, err = fmt.Fprintf(f, "R%d, [R%d,#%d]\n", list[i].rt, list[i].rn, list[i].address)
			if err != nil {
				log.Fatal(err)
			}

		case "I":
			_, err := fmt.Fprintf(f, "%s %s %s %s\t", list[i].rawInstruction[0:10],
				list[i].rawInstruction[10:22], list[i].rawInstruction[22:27],
				list[i].rawInstruction[27:32])

			_, err = fmt.Fprintf(f, "\t%d\t%s ", list[i].programCnt, list[i].op)

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

			_, err = fmt.Fprintf(f, "\t%d\t%s ", list[i].programCnt, list[i].op)

			_, err = fmt.Fprintf(f, "R%d #%d\n", list[i].conditional, list[i].offset)
			if err != nil {
				log.Fatal(err)
			}
		case "IM":
			_, err := fmt.Fprintf(f, "%s %s %s %s\t", list[i].rawInstruction[0:9],
				list[i].rawInstruction[9:11], list[i].rawInstruction[11:27],
				list[i].rawInstruction[27:32])

			_, err = fmt.Fprintf(f, "\t%d\t%s ", list[i].programCnt, list[i].op)

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
func RTypeFormat(list *Instruction) {

	// bits 12 - 16
	list.rm = uint8((list.linevalue & 2031616) >> 16)
	// bits 17 - 22
	list.shamt = uint8((list.linevalue & 64512) >> 10)
	// bits 23 - 27
	list.rn = uint8((list.linevalue & 992) >> 5)
	// bits 28 - 32
	list.rd = uint8(list.linevalue & 31)
}

func DTypeFormat(list *Instruction) {
	// bits 12 - 20
	list.address = uint16((list.linevalue & 2093056) >> 12)
	// bits 21 - 22
	list.op2 = uint8((list.linevalue & 3072) >> 10)
	// bits 23 - 27
	list.rn = uint8((list.linevalue & 992) >> 5)
	// bits 28 - 32
	list.rt = uint8(list.linevalue & 31)
}

func ITypeFormat(list *Instruction) {
	// bits 11 - 22
	list.immediate = int32(TwoComplement((list.linevalue&4193280)>>10, 12))
	// bits 23 - 27
	list.rn = uint8((list.linevalue & 992) >> 5)
	// bits 28 -32
	list.rd = uint8(list.linevalue & 31)
}

func BTypeFormat(list *Instruction) {
	// bits 7 - 32
	list.offset = int32(TwoComplement((list.linevalue & 67108863), 26))
}

func CBTypeFormat(list *Instruction) {
	// bits 9 - 27
	list.offset = int32(TwoComplement((list.linevalue&16777184)>>5, 19))
	// bits 28 - 32
	list.conditional = uint8(list.linevalue & 31)
}

func IMTypeFormat(list *Instruction) {
	// bits 10-11
	list.shiftCode = uint8((list.linevalue & 6291456) >> 21)
	// bits 12 - 27
	list.field = uint32((list.linevalue & 2097120) >> 5)
	// bits 28 - 32
	list.rd = uint8((list.linevalue & 31))
}

func Break(list *Instruction) {
	// bits 11 - 32
	list.brk = uint32(list.linevalue & 2031591)
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
			MemoryIndex[Register[list.rn]+int(list.address)*4] = int64(Register[list.rt])
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
}

func counterInSlice(c int, slice [][]int) bool {
	for i := range SturData {
		if c == SturData[i][0] {
			return true
		}
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/*
********************************************
//  main
********************************************
*/

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
var MemoryIndex = make(map[int]int64)

func main() {

	//team := "team12_out"
	//dis := "_dis.txt"
	//sim := "_sim.txt"
	//InputFileName := flag.String("i", "dtest2_bin.txt", "Gets the input file name")
	InputFileName := flag.String("i", "input.txt", "Gets the input file name")
	OutputFileName := flag.String("o", "team12_out", "Gets the output file name")
	//OutputFileName2 := flag.String("k", "team12_out", "Gets the output file name")

	flag.Parse()

	readInstruction(*InputFileName)
	processInput(InputParsed)
	initializeRegisters(InputParsed)
	ProcessSnapshot(InputParsed)
	writeInstruction(*OutputFileName+"_dis.txt", InputParsed)
	writeSimulator(*OutputFileName+"_sim.txt", InputParsed)

	fmt.Println("end project 1")
}
