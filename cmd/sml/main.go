package main

import (
	"flag"
	"fmt"
	"os"

	"shanhu.io/smlvm/arch"
	"shanhu.io/smlvm/asm"
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl"
)

func handleErrs(errs []*lexing.Error) {
	if errs == nil {
		return
	}
	for _, err := range errs {
		fmt.Println(err)
	}
	os.Exit(-1)
}

func main() {
	pkg := flag.String("pkg", "/...", "package to build")
	homeDir := flag.String("home", ".", "the home directory")
	plan := flag.Bool("plan", false, "plan only")
	golike := flag.Bool("golike", false, "uses go-like syntax")
	runTests := flag.Bool("test", true, "run tests")
	std := flag.String("std", "/smallrepo/std", "stdlib package")
	initPC := flag.Uint("initpc", arch.InitPC, "init PC register value")
	staticOnly := flag.Bool("static", false, "do static analysis only")
	flag.Parse()

	lang := pl.Lang(*golike)
	home := builds.NewDirHome(*homeDir, lang)
	home.MemHome = pl.MakeMemHome(lang)
	home.AddLang("asm", asm.Lang())

	b := builds.NewBuilder(home, home, *std)
	b.Verbose = true
	b.InitPC = uint32(*initPC)
	b.RunTests = *runTests
	b.StaticOnly = *staticOnly

	pkgs, err := builds.SelectPkgs(home, *pkg)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if !*plan {
		handleErrs(b.BuildPkgs(pkgs))
	} else {
		buildOrder, errs := b.Plan(pkgs)
		handleErrs(errs)
		for _, p := range buildOrder {
			fmt.Println(p)
		}
	}
}
