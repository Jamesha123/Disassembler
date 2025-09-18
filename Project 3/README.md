# CS3339 Project 3 - ARM Processor Disassembler and Simulator

## Team 12

A comprehensive ARM processor disassembler and simulator written in Go that can parse binary ARM instructions and simulate their execution with a pipelined processor architecture.

## Overview

This project implements a complete ARM processor simulation system that includes:
- **Disassembler**: Converts binary ARM instructions to human-readable assembly code
- **Simulator**: Executes the disassembled instructions with a pipelined processor model
- **Cache System**: Implements a 2-way set-associative cache with LRU replacement
- **Pipeline Stages**: Fetch, Issue, ALU, Memory, and Writeback stages

## Features

### Supported Instruction Types
- **R-Type**: AND, ADD, ORR, SUB, LSR, LSL, ASR, EOR
- **I-Type**: ADDI, SUBI
- **D-Type**: LDUR (load), STUR (store)
- **B-Type**: Branch instructions
- **CB-Type**: CBZ (compare and branch if zero), CBNZ (compare and branch if not zero)
- **IM-Type**: MOVZ, MOVK (move with zero/keep)
- **Special**: BREAK, NOP

### Architecture Components
- **32 General Purpose Registers** (R0-R31)
- **Program Counter** management
- **Memory System** with cache simulation
- **Pipeline Buffers** for instruction flow
- **ALU Operations** for arithmetic and logical instructions

## File Structure

```
Project 3/
├── Project Main.go          # Main entry point and instruction parsing
├── Disassembler.go         # Core disassembly logic and instruction formatting
├── Simulator.go            # Pipeline simulation and register management
├── ALU.go                  # Arithmetic Logic Unit operations
├── Cache.go                # Cache system implementation
├── Mem.go                  # Memory access operations
├── Issue.go                # Instruction issue logic
├── WriteBack.go            # Writeback stage implementation
├── fetch.go                # Instruction fetch stage
├── Helper.go               # Utility functions
├── go.mod                  # Go module configuration
├── input.txt               # Input binary instruction file
├── dtest2_bin.txt          # Test binary file
└── README.md               # This file
```

## Usage

### Prerequisites
- Go 1.19 or later

### Running the Program

1. **Compile the program:**
   ```bash
   go build -o disassembler .
   ```

2. **Run with default settings:**
   ```bash
   ./disassembler
   ```
   This will process `dtest2_bin.txt` and generate:
   - `team12_out_dis.txt` - Disassembled instructions
   - `team12_out_sim.txt` - Simulation results

3. **Run with custom output filename:**
   ```bash
   ./disassembler -o my_output
   ```

### Input Format
The program expects a text file containing 32-bit binary ARM instructions, one per line.

### Output Files
- **`*_dis.txt`**: Contains the disassembled instructions with:
  - Binary representation (grouped by instruction fields)
  - Program counter address
  - Assembly mnemonic
  - Operands (registers, immediates, addresses)
  
- **`*_sim.txt`**: Contains simulation results including:
  - Register states at each cycle
  - Program counter values
  - Memory contents
  - Cache statistics

## Instruction Format Support

### R-Type Instructions
```
Format: [11:opcode][5:rm][6:shamt][5:rn][5:rd]
Example: ADD R1, R2, R3
```

### I-Type Instructions
```
Format: [10:opcode][12:immediate][5:rn][5:rd]
Example: ADDI R1, R2, #100
```

### D-Type Instructions
```
Format: [11:opcode][9:address][2:op2][5:rn][5:rt]
Example: LDUR R1, [R2, #4]
```

### B-Type Instructions
```
Format: [6:opcode][26:offset]
Example: B #1000
```

### CB-Type Instructions
```
Format: [8:opcode][19:offset][5:conditional]
Example: CBZ R1, #50
```

## Pipeline Architecture

The simulator implements a 5-stage pipeline:

1. **Fetch**: Retrieves instructions from memory
2. **Issue**: Decodes and issues instructions to appropriate units
3. **ALU**: Performs arithmetic and logical operations
4. **Memory**: Handles load/store operations
5. **Writeback**: Updates registers with results

### Pipeline Buffers
- **Pre-Issue Buffer**: 4 instructions
- **Pre-ALU Buffer**: 2 instructions
- **Pre-Memory Buffer**: 2 instructions
- **Post-ALU Buffer**: 1 result
- **Post-Memory Buffer**: 1 result

## Cache System

- **Type**: 2-way set-associative
- **Size**: 4 sets × 2 blocks × 2 words per block
- **Replacement Policy**: LRU (Least Recently Used)
- **Block Size**: 8 bytes (2 words)

## Example Output

### Disassembled Instructions
```
00000000000 00001 000000 00010 00001    96    ADD R1, R2, R1
00000000000 00011 000000 00010 00010    100   ADD R2, R2, R3
```

### Simulation Results
```
====================
cycle:1
registers:
R00: 5    R01: 8    R02: 200  R03: 0    R04: 0    R05: 0    R06: 0    R07: 10
...
PC: 96
```

## Development

### Adding New Instructions
1. Add opcode matching in `opcodeMatching()` function
2. Implement instruction format parsing
3. Add ALU operation if needed
4. Update output formatting

### Testing
The project includes test files:
- `dtest2_bin.txt`: Sample binary instructions
- Various output files for verification

## Team Information
- **Course**: CS3339 - Computer Architecture
- **Project**: ARM Processor Disassembler and Simulator
- **Team**: 12

## License
This project is developed for educational purposes as part of CS3339 coursework.