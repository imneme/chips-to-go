// z80/pins.go
package z80

// Address pins
const (
	PIN_A0   = 0
	PIN_A1   = 1
	PIN_A2   = 2
	PIN_A3   = 3
	PIN_A4   = 4
	PIN_A5   = 5
	PIN_A6   = 6
	PIN_A7   = 7
	PIN_A8   = 8
	PIN_A9   = 9
	PIN_A10  = 10
	PIN_A11  = 11
	PIN_A12  = 12
	PIN_A13  = 13
	PIN_A14  = 14
	PIN_A15  = 15
	PIN_D0   = 16
	PIN_D1   = 17
	PIN_D2   = 18
	PIN_D3   = 19
	PIN_D4   = 20
	PIN_D5   = 21
	PIN_D6   = 22
	PIN_D7   = 23
	PIN_M1   = 24 // machine cycle 1
	PIN_MREQ = 25 // memory request
	PIN_IORQ = 26 // input/output request
	PIN_RD   = 27 // read
	PIN_WR   = 28 // write
	PIN_HALT = 29 // halt state
	PIN_INT  = 30 // interrupt request
	PIN_RES  = 31 // reset requested
	PIN_NMI  = 32 // non-maskable interrupt
	PIN_WAIT = 33 // wait requested
	PIN_RFSH = 34 // refresh
)

// Pin masks
const (
	A0   = uint64(1) << PIN_A0
	A1   = uint64(1) << PIN_A1
	A2   = uint64(1) << PIN_A2
	A3   = uint64(1) << PIN_A3
	A4   = uint64(1) << PIN_A4
	A5   = uint64(1) << PIN_A5
	A6   = uint64(1) << PIN_A6
	A7   = uint64(1) << PIN_A7
	A8   = uint64(1) << PIN_A8
	A9   = uint64(1) << PIN_A9
	A10  = uint64(1) << PIN_A10
	A11  = uint64(1) << PIN_A11
	A12  = uint64(1) << PIN_A12
	A13  = uint64(1) << PIN_A13
	A14  = uint64(1) << PIN_A14
	A15  = uint64(1) << PIN_A15
	D0   = uint64(1) << PIN_D0
	D1   = uint64(1) << PIN_D1
	D2   = uint64(1) << PIN_D2
	D3   = uint64(1) << PIN_D3
	D4   = uint64(1) << PIN_D4
	D5   = uint64(1) << PIN_D5
	D6   = uint64(1) << PIN_D6
	D7   = uint64(1) << PIN_D7
	M1   = uint64(1) << PIN_M1
	MREQ = uint64(1) << PIN_MREQ
	IORQ = uint64(1) << PIN_IORQ
	RD   = uint64(1) << PIN_RD
	WR   = uint64(1) << PIN_WR
	HALT = uint64(1) << PIN_HALT
	INT  = uint64(1) << PIN_INT
	RES  = uint64(1) << PIN_RES
	NMI  = uint64(1) << PIN_NMI
	WAIT = uint64(1) << PIN_WAIT
	RFSH = uint64(1) << PIN_RFSH

	CTRL_PIN_MASK = M1 | MREQ | IORQ | RD | WR | RFSH
	PIN_MASK      = (uint64(1) << 40) - 1
)

// Helper functions to manipulate pin states
func MakePins(ctrl, addr uint64, data uint8) uint64 {
	return ctrl | (uint64(data)&0xFF)<<16 | (addr & 0xFFFF)
}

func GetAddr(pins uint64) uint16 {
	return uint16(pins)
}

func SetAddr(pins *uint64, addr uint16) {
	*pins = (*pins &^ 0xFFFF) | uint64(addr&0xFFFF)
}

func GetData(pins uint64) uint8 {
	return uint8(pins >> 16)
}

func SetData(pins *uint64, data uint8) {
	*pins = (*pins &^ 0xFF0000) | (uint64(data) << 16 & 0xFF0000)
}
