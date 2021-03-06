package asm

import (
	"bytes"
	"math"
	"strconv"

	"shanhu.io/smlvm/arch"
	"shanhu.io/smlvm/asm/parse"
	"shanhu.io/smlvm/lexing"
)

func nbitAlign(nbit int) uint32 {
	if nbit == 8 {
		return 0
	} else if nbit == 32 {
		return 4
	} else {
		panic("invalid nbit")
	}
}

const (
	modeSigned = 0x1
	modeWord   = 0x2
	modeFloat  = 0x4
)

func parseDataNums(p lexing.Logger, args []*lexing.Token, mode int) (
	[]byte, uint32,
) {
	if !checkTypeAll(p, args, parse.Operand) {
		return nil, 0
	}

	var ui uint32
	nbit := 8
	if mode&modeWord != 0 {
		nbit = 32
	}
	var e error

	buf := new(bytes.Buffer)
	for _, arg := range args {
		if mode&modeFloat != 0 {
			var f float64
			f, e = strconv.ParseFloat(arg.Lit, 32)
			ui = math.Float32bits(float32(f))
		} else if mode&modeSigned != 0 {
			var i int64
			i, e = strconv.ParseInt(arg.Lit, 0, nbit)
			ui = uint32(i)
		} else {
			var ui64 uint64
			ui64, e = strconv.ParseUint(arg.Lit, 0, nbit)
			ui = uint32(ui64)
		}
		if e != nil {
			p.Errorf(arg.Pos, "%s", e)
			return nil, 0
		}

		if nbit == 8 {
			buf.WriteByte(byte(ui))
		} else if nbit == 32 {
			var bs [4]byte
			arch.Endian.PutUint32(bs[:], ui)
			buf.Write(bs[:])
		}
	}

	return buf.Bytes(), nbitAlign(nbit)
}
