package main

import (
	"fmt"

	"github.com/imneme/chips-to-go/z80" // Local relative import for testing
)

func main() {
	// Test code here...
	cpu, pins := z80.New()
	fmt.Printf("Z80 initialized, pins: 0x%016X\n", pins)

	// Create a small test memory
	mem := make([]byte, 256)

	// Put a couple of instructions at address 0
	mem[0] = 0x3E // LD A, 42
	mem[1] = 0x2A
	mem[2] = 0x76 // HALT

	// Run until HALT
	halted := false
	for !halted {
		pins = cpu.Tick(pins)

		// Handle memory requests
		if pins&z80.MREQ != 0 {
			addr := z80.GetAddr(pins)
			if addr < uint16(len(mem)) {
				if pins&z80.RD != 0 {
					z80.SetData(&pins, mem[addr])
				} else if pins&z80.WR != 0 {
					mem[addr] = z80.GetData(pins)
				}
			}
		}

		// Check for HALT
		if pins&z80.HALT != 0 {
			halted = true
		}

		// Print register state when instruction completes
		if cpu.OpDone() {
			fmt.Printf("A: %02X  BC: %04X  DE: %04X  HL: %04X  PC: %04X\n",
				cpu.A(), cpu.BC(), cpu.DE(), cpu.HL(), cpu.PC())
		}
	}

	fmt.Println("CPU halted!")
}
