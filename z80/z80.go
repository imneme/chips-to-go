// z80/z80.go
package z80

/*
#cgo CFLAGS: -I./include
#define CHIPS_IMPL
#include "z80.h"
*/
import "C"

// CPU represents a Z80 CPU instance
type CPU struct {
	cpu C.z80_t
}

// New creates a new Z80 CPU instance and initializes it
func New() (*CPU, uint64) {
	cpu := &CPU{}
	pins := uint64(C.z80_init(&cpu.cpu))
	return cpu, pins
}

// Reset resets the CPU to its initial state
func (c *CPU) Reset() uint64 {
	return uint64(C.z80_reset(&c.cpu))
}

// Tick advances the CPU by one clock cycle
func (c *CPU) Tick(pins uint64) uint64 {
	return uint64(C.z80_tick(&c.cpu, C.uint64_t(pins)))
}

// Prefetch forces execution to continue at the specified address
func (c *CPU) Prefetch(newPC uint16) uint64 {
	return uint64(C.z80_prefetch(&c.cpu, C.uint16_t(newPC)))
}

// OpDone returns true when a full instruction has finished executing
func (c *CPU) OpDone() bool {
	return bool(C.z80_opdone(&c.cpu))
}
