// OMSE: One More Spectrum Emulator
// (c) 2025 Melissa O'Neill
// Go Port

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/imneme/chips-to-go/z80"
	"github.com/veandco/go-sdl2/sdl"
)

// Timing constants (all in T-states)
const (
	TStatesPerLine  = 224
	TStatesPerFrame = 69888     // (64+192+56)*224
	ClockRate       = 3_500_000 // 3.5MHz
)

// Memory system - implements what's needed for display
type Memory struct {
	ram []byte
}

func NewMemory() *Memory {
	mem := &Memory{
		ram: make([]byte, 0x10000), // 64K of RAM
	}

	// Create a recognizable pattern in screen memory
	for y := uint16(0); y < 192; y++ {
		for x := uint16(0); x < 32; x++ {
			addr := 0x4000 + (y * 32) + x
			// Create diagonal stripes
			if ((x + (y / 8)) & 0x07) != 0 {
				mem.ram[addr] = 0xAA
			} else {
				mem.ram[addr] = 0x55
			}
		}
	}

	// Set attributes to alternate colors
	for y := uint16(0); y < 24; y++ {
		for x := uint16(0); x < 32; x++ {
			attrAddr := 0x5800 + (y * 32) + x
			// Alternate between cyan on black and yellow on blue
			if ((x + y) & 1) != 0 {
				mem.ram[attrAddr] = 0x45
			} else {
				mem.ram[attrAddr] = 0x16
			}
		}
	}

	return mem
}

func (m *Memory) Read(address uint16) byte {
	return m.ram[address]
}

func (m *Memory) Write(address uint16, value byte) {
	if address < 0x4000 {
		return // Ignore writes to ROM
	}
	m.ram[address] = value
}

func (m *Memory) LoadFromFile(filename string, addr uint16, size uint16) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %s: %v", filename, err)
	}
	defer file.Close()

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info: %v", err)
	}

	// Check if we have enough data
	if fileInfo.Size() < int64(size) {
		return fmt.Errorf("file too small: need at least %d bytes", size)
	}

	return m.LoadFromReader(file, addr, size)
}

func (m *Memory) LoadFromReader(r io.Reader, addr uint16, size uint16) error {
	_, err := io.ReadFull(r, m.ram[addr:uint32(addr)+uint32(size)])
	return err
}

// CRT display using SDL
type CRT struct {
	window        *sdl.Window
	renderer      *sdl.Renderer
	screenTexture *sdl.Texture
	pixels        []uint32
	oddField      bool
	flashInverted bool
}

// CRT constants
const (
	TotalWidth     = 352            // 352 pixels
	Columns        = TotalWidth / 8 // 352/8 columns
	FieldLines     = 312            // Total PAL lines per field
	TopBlanking    = 16             // Lines before visible area
	BottomBlanking = 4              // Lines after visible area
	VisibleLines   = FieldLines - TopBlanking - BottomBlanking
	CRTLines       = VisibleLines * 2 // Two fields
)

// RGBA color table
var rgbaColorTable = []uint32{
	0x00000000, // Black
	0x0000FFFF, // Blue
	0xFF000000, // Red
	0xFF00FFFF, // Magenta
	0x00FF0000, // Green
	0x00FFFFFF, // Cyan
	0xFFFF0000, // Yellow
	0xFFFFFFFF, // White
}

func NewCRT() (*CRT, error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, fmt.Errorf("SDL initialization failed: %v", err)
	}

	// Scale up 2x for better visibility
	window, err := sdl.CreateWindow(
		"OMSE — One More Spectrum Emulator (Go Port)",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		TotalWidth*2,
		CRTLines,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return nil, fmt.Errorf("window creation failed: %v", err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		window.Destroy()
		return nil, fmt.Errorf("renderer creation failed: %v", err)
	}

	screenTexture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STREAMING,
		TotalWidth,
		CRTLines,
	)
	if err != nil {
		renderer.Destroy()
		window.Destroy()
		return nil, fmt.Errorf("texture creation failed: %v", err)
	}

	// Set up 2x scaling horizontally, 1x vertically
	renderer.SetScale(2.0, 1.0)

	return &CRT{
		window:        window,
		renderer:      renderer,
		screenTexture: screenTexture,
		pixels:        make([]uint32, TotalWidth*CRTLines),
		oddField:      false,
		flashInverted: false,
	}, nil
}

func (c *CRT) Close() {
	if c.screenTexture != nil {
		c.screenTexture.Destroy()
	}
	if c.renderer != nil {
		c.renderer.Destroy()
	}
	if c.window != nil {
		c.window.Destroy()
	}
	sdl.Quit()
}

// UpdatePixels updates a group of 8 pixels at the specified location
func (c *CRT) UpdatePixels(line uint32, column uint32, displayByte byte, attrByte byte) {
	// Assertions/bounds checking
	if line >= FieldLines || column >= Columns {
		return
	}

	// Ignore updates in blanking intervals
	if line < TopBlanking || line >= (TopBlanking+VisibleLines) {
		return
	}
	line -= TopBlanking // Adjust for top blanking

	// Interlace fields
	var interlacedLine uint32
	if c.oddField {
		interlacedLine = line*2 + 1
	} else {
		interlacedLine = line * 2
	}

	// Update 8 pixels at once
	offset := (interlacedLine * TotalWidth) + (column * 8)

	// To bleed into the other line
	var bleedOffset uint32
	if c.oddField {
		bleedOffset = offset - TotalWidth
	} else {
		bleedOffset = offset + TotalWidth
	}

	// Convert attribute byte
	flash := (attrByte & 0x80) != 0
	bright := (attrByte & 0x40) != 0
	paper := (attrByte >> 3) & 0x07
	ink := attrByte & 0x07

	if flash && c.flashInverted {
		// Swap paper and ink
		paper, ink = ink, paper
	}

	// Create RGB colors
	paperColor := rgbaColorTable[paper]
	inkColor := rgbaColorTable[ink]

	// Update 8 pixels - note MSB is leftmost pixel
	for bit := 7; bit >= 0; bit-- {
		pixelSet := (displayByte & (1 << bit)) != 0
		var color uint32
		if pixelSet {
			color = inkColor
		} else {
			color = paperColor
		}

		// For main scanline, set the pixel with phosphor fade from previous
		c.pixels[offset+(7-uint32(bit))] = ((c.pixels[offset+(7-uint32(bit))] >> 2) & 0x3f3f3f3f) | color

		// Other scanline is less bright, but bleeds through. Bright colors
		// bleed through more to create the bright effect.
		if !bright {
			color = ((color >> 1) & 0x7f7f7f7f) | 0xff // 50% brightness
		} else {
			color = ((color>>3)&0x07070707)*27 | 0xff // 84% brightness
		}
		c.pixels[bleedOffset+(7-uint32(bit))] = ((c.pixels[bleedOffset+(7-uint32(bit))] >> 2) & 0x3f3f3f3f) | color
	}
}

func (c *CRT) Refresh() {
	c.screenTexture.Update(nil, unsafe.Pointer(&c.pixels[0]), TotalWidth*4)
	c.renderer.Clear()
	c.renderer.Copy(c.screenTexture, nil, nil)
	c.renderer.Present()
}

func (c *CRT) ToggleFlash() {
	c.flashInverted = !c.flashInverted
}

// IODevice interface for devices that support read and write
type IODevice interface {
	Read(addr uint16) byte
	Write(addr uint16, value byte)
}

// IODeviceBus for I/O devices
type IODeviceBus struct {
	devices map[uint16]IODevice // mask -> device
}

func NewIODeviceBus() *IODeviceBus {
	return &IODeviceBus{
		devices: make(map[uint16]IODevice),
	}
}

func (b *IODeviceBus) AddDevice(mask uint16, device IODevice) {
	b.devices[mask] = device
}

func (b *IODeviceBus) Read(addr uint16) byte {
	for mask, device := range b.devices {
		if ((^addr) & mask) == mask {
			return device.Read(addr)
		}
	}
	return 0xff // Default to all bits set
}

func (b *IODeviceBus) Write(addr uint16, value byte) {
	for mask, device := range b.devices {
		if ((^addr) & mask) == mask {
			device.Write(addr, value)
			return
		}
	}
}

// CPU implementation using our Z80 wrapper
type CPU struct {
	*z80.CPU      // Embed our Z80 CPU implementation
	memory        *Memory
	bus           *IODeviceBus
	interruptFlag bool
	pins          uint64
}

func NewCPU(memory *Memory, bus *IODeviceBus) *CPU {
	z80cpu, pins := z80.New()
	return &CPU{
		CPU:    z80cpu,
		memory: memory,
		bus:    bus,
		pins:   pins, // Store initial pin state
	}
}

func (c *CPU) Tick() {
	// Update pin state with any pending interrupt
	if c.interruptFlag {
		c.pins |= z80.INT
	} else {
		c.pins &= ^z80.INT
	}

	// Perform one Z80 tick
	c.pins = c.CPU.Tick(c.pins)

	// Process memory and I/O transactions
	c.transact()
}

func (c *CPU) transact() {
	// Handle memory access
	if c.pins&z80.MREQ != 0 {
		addr := z80.GetAddr(c.pins)
		if c.pins&z80.RD != 0 {
			// Memory read
			data := c.memory.Read(addr)
			z80.SetData(&c.pins, data)
		} else if c.pins&z80.WR != 0 {
			// Memory write
			data := z80.GetData(c.pins)
			c.memory.Write(addr, data)
		}
	} else if c.pins&z80.IORQ != 0 {
		// I/O access
		if c.pins&z80.M1 != 0 {
			// Interrupt acknowledge
			z80.SetData(&c.pins, 0xFF)
		} else {
			addr := z80.GetAddr(c.pins)
			if c.pins&z80.RD != 0 {
				// IO read
				data := c.bus.Read(addr)
				z80.SetData(&c.pins, data)
			} else if c.pins&z80.WR != 0 {
				// IO write
				data := z80.GetData(c.pins)
				c.bus.Write(addr, data)
			}
		}
	}
}

func (c *CPU) SetInterrupt(status bool) {
	c.interruptFlag = status
}

func (c *CPU) SetPC(addr uint16) {
	c.pins = c.CPU.Prefetch(addr)
}

// ULA (Uncommitted Logic Array) - the Spectrum's custom chip
type ULA struct {
	memory       *Memory
	crt          *CRT
	cpu          *CPU
	borderColor  byte
	flashFlipper byte

	// Current position tracking
	line          uint32 // Current scanline (0-311)
	lineCycle     uint32 // Current cycle within line (0-223)
	currentColumn uint32
}

const (
	ScreenStartLine    = 64
	ScreenStartColumn  = 6 // 48 pixels / 8
	ScreenWidthBytes   = 32
	ScreenHeight       = 192
	BorderTStates      = ScreenStartColumn * 4
	ScreenWidthTStates = ScreenWidthBytes * 4
	FlashRate          = 16
	InterruptDuration  = 32
)

func NewULA(memory *Memory, cpu *CPU, crt *CRT) *ULA {
	return &ULA{
		memory:        memory,
		crt:           crt,
		cpu:           cpu,
		borderColor:   0,
		flashFlipper:  FlashRate,
		line:          0,
		lineCycle:     BorderTStates,
		currentColumn: 0,
	}
}

// Read and write to I/O ports
func (u *ULA) Read(addr uint16) byte {
	return 0xff
}

func (u *ULA) Write(addr uint16, value byte) {
	u.SetBorderColor(value)
}

func (u *ULA) Tick() {
	u.cpu.Tick()

	// Check if we're in the visible (non-blanking) area
	visible := (u.line >= TopBlanking) && (u.line < (FieldLines - BottomBlanking))

	if visible && u.lineCycle < (Columns*4) {
		// Every 4 cycles we output 8 pixels
		inScreenLine := (u.line >= ScreenStartLine) && (u.line < (ScreenStartLine + ScreenHeight))

		if u.lineCycle%4 == 0 {
			u.currentColumn = u.lineCycle / 4
			inScreenCol := (u.currentColumn >= ScreenStartColumn) &&
				(u.currentColumn < (ScreenStartColumn + ScreenWidthBytes))

			if inScreenLine && inScreenCol {
				// We're in the actual screen area - fetch and display pixels
				screenLine := u.line - ScreenStartLine
				screenCol := u.currentColumn - ScreenStartColumn

				addr := u.calculateDisplayAddress(screenLine, screenCol)
				displayByte := u.memory.Read(addr)
				attrByte := u.memory.Read(u.calculateAttrAddress(screenLine, screenCol))

				u.crt.UpdatePixels(u.line, u.currentColumn, displayByte, attrByte)
			} else {
				// Border area
				borderAttr := (u.borderColor << 3)
				u.crt.UpdatePixels(u.line, u.currentColumn, 0x00, borderAttr)
			}
		}
	}

	// Update position counters
	u.lineCycle++
	if u.line == 0 && u.lineCycle == BorderTStates {
		u.cpu.SetInterrupt(true)
	} else if u.line == 0 && u.lineCycle == BorderTStates+InterruptDuration {
		u.cpu.SetInterrupt(false)
	}

	if u.lineCycle >= TStatesPerLine {
		u.lineCycle = 0
		u.line++
		if u.line >= FieldLines {
			u.line = 0
			u.flashFlipper--
			if u.flashFlipper == 0 {
				u.flashFlipper = FlashRate
				u.crt.ToggleFlash()
			}
		}
	}
}

func (u *ULA) SetBorderColor(color byte) {
	u.borderColor = color & 0x07
}

func (u *ULA) GetBorderColor() byte {
	return u.borderColor
}

func (u *ULA) calculateDisplayAddress(line, col uint32) uint16 {
	// Start of screen memory
	addr := uint16(0x4000)

	// Add Y portion
	addr |= uint16((line & 0xC0) << 5) // Which third of the screen
	addr |= uint16((line & 0x07) << 8) // Which character cell row
	addr |= uint16((line & 0x38) << 2) // Remaining bits wherever

	// Add X portion
	addr |= uint16(col & 0b00011111) // 5 bits of X go to bits 0-4

	return addr
}

func (u *ULA) calculateAttrAddress(line, col uint32) uint16 {
	return 0x5800 + uint16((line>>3)*32) + uint16(col)
}

// System combines all components
type System struct {
	memory        *Memory
	bus           *IODeviceBus
	crt           *CRT
	cpu           *CPU
	ula           *ULA
	currentTState uint64
}

const (
	ChunkSize = 13 * 8 * 224 // Execute this many T-states at once
)

func NewSystem() (*System, error) {
	memory := NewMemory()
	bus := NewIODeviceBus()

	crt, err := NewCRT()
	if err != nil {
		return nil, err
	}

	cpu := NewCPU(memory, bus)
	ula := NewULA(memory, cpu, crt)

	// Initialize subsystems
	bus.AddDevice(0x0001, ula)

	return &System{
		memory:        memory,
		bus:           bus,
		crt:           crt,
		cpu:           cpu,
		ula:           ula,
		currentTState: 0,
	}, nil
}

func (s *System) Close() {
	s.crt.Close()
}

func (s *System) Run() error {
	quit := false

	// Track both virtual and real time
	startTime := time.Now()
	nextRefreshTState := s.currentTState

	for !quit {
		// Handle SDL events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				quit = true
			}
		}

		// Process a chunk of cycles
		targetTState := s.currentTState + ChunkSize
		for s.currentTState < targetTState {
			s.ula.Tick()
			s.currentTState++
		}

		// Check if we need to refresh the display
		if s.currentTState >= nextRefreshTState {
			s.crt.Refresh()
			nextRefreshTState += TStatesPerFrame
		}

		// Sleep if we're ahead
		elapsedTime := time.Since(startTime)
		expectedTime := time.Duration(s.currentTState*1000000/ClockRate) * time.Microsecond

		if expectedTime > elapsedTime {
			aheadBy := expectedTime - elapsedTime
			if aheadBy.Milliseconds() > 1 {
				time.Sleep(aheadBy - time.Millisecond)
			}
		}
	}

	return nil
}

func (s *System) LoadSNA(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %s: %v", filename, err)
	}
	defer file.Close()

	readByte := func() (byte, error) {
		var b [1]byte
		_, err := file.Read(b[:])
		return b[0], err
	}

	readWord := func() (uint16, error) {
		var b [2]byte
		_, err := file.Read(b[:])
		return uint16(b[0]) | (uint16(b[1]) << 8), err
	}

	// Set PC to the standard SNA return address
	s.cpu.SetPC(0x0072)

	// Read registers
	var err2 error

	// I register
	i, err := readByte()
	if err != nil {
		return err
	}
	// We need to add a method to set the I register in our Z80 wrapper
	s.cpu.SetI(i)

	// Alternative register set
	hl2, err2 := readWord()
	if err2 != nil {
		return err2
	}
	s.cpu.SetHL2(hl2)

	de2, err2 := readWord()
	if err2 != nil {
		return err2
	}
	s.cpu.SetDE2(de2)

	bc2, err2 := readWord()
	if err2 != nil {
		return err2
	}
	s.cpu.SetBC2(bc2)

	af2, err2 := readWord()
	if err2 != nil {
		return err2
	}
	s.cpu.SetAF2(af2)

	// Main register set
	hl, err2 := readWord()
	if err2 != nil {
		return err2
	}
	s.cpu.SetHL(hl)

	de, err2 := readWord()
	if err2 != nil {
		return err2
	}
	s.cpu.SetDE(de)

	bc, err2 := readWord()
	if err2 != nil {
		return err2
	}
	s.cpu.SetBC(bc)

	iy, err2 := readWord()
	if err2 != nil {
		return err2
	}
	s.cpu.SetIY(iy)

	ix, err2 := readWord()
	if err2 != nil {
		return err2
	}
	s.cpu.SetIX(ix)

	// Interrupt status
	intByte, err := readByte()
	if err != nil {
		return err
	}
	s.cpu.SetIFF2((intByte & 0x04) != 0)

	// R register
	r, err := readByte()
	if err != nil {
		return err
	}
	s.cpu.SetR(r)

	// AF and SP registers
	af, err2 := readWord()
	if err2 != nil {
		return err2
	}
	s.cpu.SetAF(af)

	sp, err2 := readWord()
	if err2 != nil {
		return err2
	}
	s.cpu.SetSP(sp)

	// Interrupt mode
	im, err := readByte()
	if err != nil {
		return err
	}
	s.cpu.SetIM(im)

	// Border color
	borderColor, err := readByte()
	if err != nil {
		return err
	}
	s.ula.SetBorderColor(borderColor)

	// Load RAM
	return s.memory.LoadFromReader(file, 0x4000, 49152)
}

func main() {
	// Parse command line arguments
	romLoaded := false

	system, err := NewSystem()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer system.Close()

	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			arg := os.Args[i]

			if arg == "-h" || arg == "--help" {
				fmt.Printf("Usage: %s [options] [filename...]\n"+
					"Options:\n"+
					"  -h, --help           Show this help message\n"+
					"If no filename is provided, boot into 48.rom.\n\n"+
					"(.scr, .rom and .sna files are supported)\n", os.Args[0])
				return
			} else if filepath.Ext(arg) == ".rom" {
				// Load the ROM file into memory
				err := system.memory.LoadFromFile(arg, 0x0000, 16384)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
				romLoaded = true
			} else if filepath.Ext(arg) == ".sna" {
				// Load the SNA file into memory
				err := system.LoadSNA(arg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
			} else if filepath.Ext(arg) == ".scr" {
				// Load the SCR file into memory
				err := system.memory.LoadFromFile(arg, 0x4000, 6912)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Unknown file type: %s\n", arg)
				os.Exit(1)
			}
		}
	}

	// Load the ROM (48.rom) into memory if not already loaded
	if !romLoaded {
		err := system.memory.LoadFromFile("48.rom", 0x0000, 16384)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			// Continue without ROM, we'll just use the pattern
		}
	}

	err = system.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
