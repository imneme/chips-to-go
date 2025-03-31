# CHIPS to GO

This repository provides a binding of the Andre Weissflog (Floooh)'s CHIPS
Z80 core to the Go programming language. The core is a pure Go implementation.

To ensure you have the most recent Z80 core, you can use the following command to download it:

```bash
cd include
curl -OL https://raw.githubusercontent.com/floooh/chips/refs/heads/master/chips/z80.h
cd ..
```

and then experiment with the examples in the `examples` directory. To run the
emulator example, you need to do:

```
cd examples
curl -o 48.rom -L https://github.com/spectrumforeveryone/zx-roms/raw/refs/heads/main/spectrum16-48/spec48.rom
```

and then you can run

```
go run omse-mini.go
```

This will run a Go port of “One More Spectrum Emulator” (OMSE) which is a bare-bones ZX Spectrum emulator.


