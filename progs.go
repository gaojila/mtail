package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"
)

func LoadProgs(progs string) (*engine, int) {
	fis, err := ioutil.ReadDir(progs)
	if err != nil {
		log.Fatalf("Failed to list programs in %q: %s", progs, err)
	}

	e := &engine{}
	errors := 0
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if filepath.Ext(fi.Name()) != ".em" {
			continue
		}
		f, err := os.Open(fmt.Sprintf("%s/%s", progs, fi.Name()))
		if err != nil {
			log.Printf("Failed to read program %q: %s\n", fi.Name(), err)
			continue
		}
		defer f.Close()
		v, errs := Compile(fi.Name(), f)
		if errs != nil {
			errors = 1
			for _, e := range errs {
				log.Print(e)
			}
			continue
		}
		if *dump_bytecode {
			fmt.Printf("Prog %s\n", fi.Name())
			fmt.Println("Metrics")
			for i, m := range metrics {
				if m.Program == v.name {
					fmt.Printf(" %8d %s\n", i, m)
				}
			}
			fmt.Println("REs")
			for i, re := range v.re {
				fmt.Printf(" %8d /%s/\n", i, re)
			}
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)

			fmt.Fprintln(w, "disasm\tl\top\topnd\t")
			for n, i := range v.prog {
				fmt.Fprintf(w, "\t%d\t%s\t%d\t\n", n, opNames[i.op], i.opnd)
			}
			w.Flush()
		}
		e.addVm(v)
		log.Printf("loaded %s", fi.Name())
	}
	return e, errors
}