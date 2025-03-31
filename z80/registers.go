// z80/registers.go
package z80

/*
#include "z80.h"

// Helper C functions to access Z80 struct fields
static uint16_t z80_get_pc(z80_t* cpu) { return cpu->pc; }
static uint16_t z80_get_af(z80_t* cpu) { return cpu->af; }
static uint16_t z80_get_bc(z80_t* cpu) { return cpu->bc; }
static uint16_t z80_get_de(z80_t* cpu) { return cpu->de; }
static uint16_t z80_get_hl(z80_t* cpu) { return cpu->hl; }
static uint16_t z80_get_sp(z80_t* cpu) { return cpu->sp; }
static uint16_t z80_get_ix(z80_t* cpu) { return cpu->ix; }
static uint16_t z80_get_iy(z80_t* cpu) { return cpu->iy; }
static uint16_t z80_get_af2(z80_t* cpu) { return cpu->af2; }
static uint16_t z80_get_bc2(z80_t* cpu) { return cpu->bc2; }
static uint16_t z80_get_de2(z80_t* cpu) { return cpu->de2; }
static uint16_t z80_get_hl2(z80_t* cpu) { return cpu->hl2; }

static uint8_t z80_get_a(z80_t* cpu) { return cpu->a; }
static uint8_t z80_get_f(z80_t* cpu) { return cpu->f; }
static uint8_t z80_get_b(z80_t* cpu) { return cpu->b; }
static uint8_t z80_get_c(z80_t* cpu) { return cpu->c; }
static uint8_t z80_get_d(z80_t* cpu) { return cpu->d; }
static uint8_t z80_get_e(z80_t* cpu) { return cpu->e; }
static uint8_t z80_get_h(z80_t* cpu) { return cpu->h; }
static uint8_t z80_get_l(z80_t* cpu) { return cpu->l; }
static uint8_t z80_get_i(z80_t* cpu) { return cpu->i; }
static uint8_t z80_get_r(z80_t* cpu) { return cpu->r; }
static bool z80_get_iff1(z80_t* cpu) { return cpu->iff1; }
static bool z80_get_iff2(z80_t* cpu) { return cpu->iff2; }
static uint8_t z80_get_im(z80_t* cpu) { return cpu->im; }

// Setters

static void z80_set_pc(z80_t* cpu, uint16_t pc) { cpu->pc = pc; }
static void z80_set_af(z80_t* cpu, uint16_t af) { cpu->af = af; }
static void z80_set_bc(z80_t* cpu, uint16_t bc) { cpu->bc = bc; }
static void z80_set_de(z80_t* cpu, uint16_t de) { cpu->de = de; }
static void z80_set_hl(z80_t* cpu, uint16_t hl) { cpu->hl = hl; }
static void z80_set_sp(z80_t* cpu, uint16_t sp) { cpu->sp = sp; }
static void z80_set_ix(z80_t* cpu, uint16_t ix) { cpu->ix = ix; }
static void z80_set_iy(z80_t* cpu, uint16_t iy) { cpu->iy = iy; }
static void z80_set_af2(z80_t* cpu, uint16_t af2) { cpu->af2 = af2; }
static void z80_set_bc2(z80_t* cpu, uint16_t bc2) { cpu->bc2 = bc2; }
static void z80_set_de2(z80_t* cpu, uint16_t de2) { cpu->de2 = de2; }
static void z80_set_hl2(z80_t* cpu, uint16_t hl2) { cpu->hl2 = hl2; }

static void z80_set_a(z80_t* cpu, uint8_t a) { cpu->a = a; }
static void z80_set_f(z80_t* cpu, uint8_t f) { cpu->f = f; }
static void z80_set_b(z80_t* cpu, uint8_t b) { cpu->b = b; }
static void z80_set_c(z80_t* cpu, uint8_t c) { cpu->c = c; }
static void z80_set_d(z80_t* cpu, uint8_t d) { cpu->d = d; }
static void z80_set_e(z80_t* cpu, uint8_t e) { cpu->e = e; }
static void z80_set_h(z80_t* cpu, uint8_t h) { cpu->h = h; }
static void z80_set_l(z80_t* cpu, uint8_t l) { cpu->l = l; }
static void z80_set_i(z80_t* cpu, uint8_t i) { cpu->i = i; }
static void z80_set_r(z80_t* cpu, uint8_t r) { cpu->r = r; }
static void z80_set_iff1(z80_t* cpu, bool iff1) { cpu->iff1 = iff1; }
static void z80_set_iff2(z80_t* cpu, bool iff2) { cpu->iff2 = iff2; }
static void z80_set_im(z80_t* cpu, uint8_t im) { cpu->im = im; }


*/
import "C"

// GetPC returns the program counter
func (c *CPU) PC() uint16 {
	return uint16(C.z80_get_pc(&c.cpu))
}

// GetAF returns the combined A and F register
func (c *CPU) AF() uint16  { return uint16(C.z80_get_af(&c.cpu)) }
func (c *CPU) BC() uint16  { return uint16(C.z80_get_bc(&c.cpu)) }
func (c *CPU) DE() uint16  { return uint16(C.z80_get_de(&c.cpu)) }
func (c *CPU) HL() uint16  { return uint16(C.z80_get_hl(&c.cpu)) }
func (c *CPU) SP() uint16  { return uint16(C.z80_get_sp(&c.cpu)) }
func (c *CPU) IX() uint16  { return uint16(C.z80_get_ix(&c.cpu)) }
func (c *CPU) IY() uint16  { return uint16(C.z80_get_iy(&c.cpu)) }
func (c *CPU) AF2() uint16 { return uint16(C.z80_get_af2(&c.cpu)) }
func (c *CPU) BC2() uint16 { return uint16(C.z80_get_bc2(&c.cpu)) }
func (c *CPU) DE2() uint16 { return uint16(C.z80_get_de2(&c.cpu)) }
func (c *CPU) HL2() uint16 { return uint16(C.z80_get_hl2(&c.cpu)) }

// Individual register access
func (c *CPU) A() uint8   { return uint8(C.z80_get_a(&c.cpu)) }
func (c *CPU) F() uint8   { return uint8(C.z80_get_f(&c.cpu)) }
func (c *CPU) B() uint8   { return uint8(C.z80_get_b(&c.cpu)) }
func (c *CPU) C() uint8   { return uint8(C.z80_get_c(&c.cpu)) }
func (c *CPU) D() uint8   { return uint8(C.z80_get_d(&c.cpu)) }
func (c *CPU) E() uint8   { return uint8(C.z80_get_e(&c.cpu)) }
func (c *CPU) H() uint8   { return uint8(C.z80_get_h(&c.cpu)) }
func (c *CPU) L() uint8   { return uint8(C.z80_get_l(&c.cpu)) }
func (c *CPU) I() uint8   { return uint8(C.z80_get_i(&c.cpu)) }
func (c *CPU) R() uint8   { return uint8(C.z80_get_r(&c.cpu)) }
func (c *CPU) IFF1() bool { return bool(C.z80_get_iff1(&c.cpu)) }
func (c *CPU) IFF2() bool { return bool(C.z80_get_iff2(&c.cpu)) }
func (c *CPU) IM() uint8  { return uint8(C.z80_get_im(&c.cpu)) }

// Setters

func (c *CPU) SetPC(pc uint16)   { C.z80_set_pc(&c.cpu, C.uint16_t(pc)) }
func (c *CPU) SetAF(af uint16)   { C.z80_set_af(&c.cpu, C.uint16_t(af)) }
func (c *CPU) SetBC(bc uint16)   { C.z80_set_bc(&c.cpu, C.uint16_t(bc)) }
func (c *CPU) SetDE(de uint16)   { C.z80_set_de(&c.cpu, C.uint16_t(de)) }
func (c *CPU) SetHL(hl uint16)   { C.z80_set_hl(&c.cpu, C.uint16_t(hl)) }
func (c *CPU) SetSP(sp uint16)   { C.z80_set_sp(&c.cpu, C.uint16_t(sp)) }
func (c *CPU) SetIX(ix uint16)   { C.z80_set_ix(&c.cpu, C.uint16_t(ix)) }
func (c *CPU) SetIY(iy uint16)   { C.z80_set_iy(&c.cpu, C.uint16_t(iy)) }
func (c *CPU) SetAF2(af2 uint16) { C.z80_set_af2(&c.cpu, C.uint16_t(af2)) }
func (c *CPU) SetBC2(bc2 uint16) { C.z80_set_bc2(&c.cpu, C.uint16_t(bc2)) }
func (c *CPU) SetDE2(de2 uint16) { C.z80_set_de2(&c.cpu, C.uint16_t(de2)) }
func (c *CPU) SetHL2(hl2 uint16) { C.z80_set_hl2(&c.cpu, C.uint16_t(hl2)) }

// Individual register access
func (c *CPU) SetA(a uint8)  { C.z80_set_a(&c.cpu, C.uint8_t(a)) }
func (c *CPU) SetF(f uint8)  { C.z80_set_f(&c.cpu, C.uint8_t(f)) }
func (c *CPU) SetB(b uint8)  { C.z80_set_b(&c.cpu, C.uint8_t(b)) }
func (c *CPU) SetC(cr uint8) { C.z80_set_c(&c.cpu, C.uint8_t(cr)) }
func (c *CPU) SetD(d uint8)  { C.z80_set_d(&c.cpu, C.uint8_t(d)) }
func (c *CPU) SetE(e uint8)  { C.z80_set_e(&c.cpu, C.uint8_t(e)) }
func (c *CPU) SetH(h uint8)  { C.z80_set_h(&c.cpu, C.uint8_t(h)) }
func (c *CPU) SetL(l uint8)  { C.z80_set_l(&c.cpu, C.uint8_t(l)) }
func (c *CPU) SetI(i uint8)  { C.z80_set_i(&c.cpu, C.uint8_t(i)) }
func (c *CPU) SetR(r uint8)  { C.z80_set_r(&c.cpu, C.uint8_t(r)) }

func (c *CPU) SetIFF1(iff1 bool) { C.z80_set_iff1(&c.cpu, C._Bool(iff1)) }
func (c *CPU) SetIFF2(iff2 bool) { C.z80_set_iff2(&c.cpu, C._Bool(iff2)) }
func (c *CPU) SetIM(im uint8)    { C.z80_set_im(&c.cpu, C.uint8_t(im)) }
