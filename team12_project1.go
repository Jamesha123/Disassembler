package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	//"testing/quick"
)

type Instruction struct {
	typeofInstruction string
	rawInstruction    string
	linevalue         uint64
	programCnt        int
	opcode            uint64
	op                string
	op2               uint8
	rd                uint8
	rn                uint8
	rm                uint8
	rt                uint8
	shamt             uint8
	im                string
	immediate         int16
	offset            int32
	conditional       uint8
	address           uint16
	shiftCode         uint8
	field             uint32
	brk               uint32
	bitValue          int64
}

/*
*************************************
// functions
*************************************
*/
func processInput(list []Instruction) {

	for i := 0; i < len(list); i++ {
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
		case "NUM":
			var val uint64
			val, _ = strconv.ParseUint(list[i].rawInstruction, 2, 32)
			list[i].bitValue = TwoComplement(val, 32)
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
	} else if ins.opcode == 0 {
		ins.op = "NOP"
		ins.typeofInstruction = "NOP"
	} else if ins.opcode == 1872 {
		ins.op = "EOR"
		ins.typeofInstruction = "R"
	} else if ins.opcode == 2038 {
		ins.op = "Break"
		ins.typeofInstruction = "BREAK"
	} else {
		ins.typeofInstruction = "NUM"
		fmt.Println("Invalid opcode")
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

			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].programCnt, list[i].op)

			_, err = fmt.Fprintf(f, "R%d, R%d, #%d\n", list[i].rd, list[i].rn, list[i].immediate)
			if err != nil {
				log.Fatal(err)
			}
		case "B":
			_, err := fmt.Fprintf(f, "%s %s\t", list[i].rawInstruction[0:6],
				list[i].rawInstruction[6:32])

			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].programCnt, list[i].op)

			_, err = fmt.Fprintf(f, "#%d\n", list[i].offset)
			if err != nil {
				log.Fatal(err)
			}
		case "CB":
			_, err := fmt.Fprintf(f, "%s %s %s\t", list[i].rawInstruction[0:8],
				list[i].rawInstruction[8:27], list[i].rawInstruction[27:32])

			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].programCnt, list[i].op)

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
			_, err = fmt.Fprintf(f, "%s %d Break\n", list[i].rawInstruction, list[i].programCnt)
			if err != nil {
				log.Fatal(err)
			}
		case "NOP":
			_, err = fmt.Fprintf(f, "%s %d No Instruction\n", list[i].rawInstruction, list[i].programCnt)
			if err != nil {
				log.Fatal(err)
			}
		case "NUM":
			_, err = fmt.Fprintf(f, "%s %d %d\n", list[i].rawInstruction, list[i].programCnt, list[i].bitValue)
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
	ins.immediate = int16(TwoComplement((ins.linevalue&4193280)>>10, 12))
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

/*
********************************************
//  main
********************************************
*/

var InputParsed []Instruction
var Breaknow bool = false
var InputFileName *string
var OutputFileName *string

func main() {

	InputFileName := flag.String("i", "dtest2_bin.txt", "Gets the input file name")
	OutputFileName := flag.String("o", "team12_addtest1_out.dis.txt", "Gets the output file name")

	flag.Parse()

	//os.Args[0] = "addtest1_bin.txt"
	//os.Args[0] = "dtest2_bin.txt"

	readInstruction(*InputFileName)
	processInput(InputParsed)
	writeInstruction(*OutputFileName, InputParsed)

}
