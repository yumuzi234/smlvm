package sempass

import (
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func unaryOpConst(b *Builder, opTok *lex8.Token, B tast.Expr) tast.Expr {
	op := opTok.Lit
	bref := tast.ExprRef(B)
	if bref.List != nil {
		b.Errorf(opTok.Pos, "invalid operation: %q on %s", op, bref.T)
		return nil
	}

	v, ok := types.NumConst(bref.T)
	if !ok {
		// TODO: support type const
		b.Errorf(opTok.Pos, "typed const operation not implemented")
		return nil
	}

	switch op {
	case "+":
		return B // just shortcut this
	case "-":
		return &tast.Const{tast.NewRef(types.NewNumber(-v))}
	}

	b.Errorf(opTok.Pos, "invalid operation: %q on %s", op, B)
	return nil
}

func binaryOpConst(b *Builder, opTok *lex8.Token, A, B tast.Expr) tast.Expr {
	op := opTok.Lit
	aref := tast.ExprRef(A)
	bref := tast.ExprRef(B)
	if aref.List != nil || bref.List != nil {
		b.Errorf(opTok.Pos, "invalid %s %q %s", aref.T, op, bref.T)
		return nil
	}

	va, oka := types.NumConst(aref.T)
	vb, okb := types.NumConst(bref.T)
	if !(oka && okb) {
		b.Errorf(opTok.Pos, "non-numeric consts ops not implemented")
		return nil
	}

	r := func(v int64) tast.Expr {
		return &tast.Const{tast.NewRef(types.NewNumber(v))}
	}

	switch op {
	case "+":
		return r(va + vb)
	case "-":
		return r(va - vb)
	case "*":
		return r(va * vb)
	case "&":
		return r(va & vb)
	case "|":
		return r(va | vb)
	case "^":
		return r(va ^ vb)
	case "%":
		if vb == 0 {
			b.Errorf(opTok.Pos, "modular by zero")
			return nil
		}
		return r(va % vb)
	case "/":
		if vb == 0 {
			b.Errorf(opTok.Pos, "divide by zero")
			return nil
		}
		return r(va / vb)
	case "<<":
		if vb < 0 {
			b.Errorf(opTok.Pos, "shift with negative value", vb)
			return nil
		}
		return r(va << uint64(vb))
	case ">>":
		if vb < 0 {
			b.Errorf(opTok.Pos, "shift with negative value", vb)
			return nil
		}
		return r(va >> uint64(vb))
	}

	b.Errorf(opTok.Pos, "%q on consts", op)
	return nil
}

func unaryOpInt(b *Builder, opTok *lex8.Token, B tast.Expr) tast.Expr {
	op := opTok.Lit
	switch op {
	case "+":
		return B
	case "-", "^":
		t := tast.ExprRef(B).T
		return &tast.OpExpr{nil, opTok, B, tast.NewRef(t)}
	}

	b.Errorf(opTok.Pos, "invalid operation: %q on %s", op, B)
	return nil
}
