package main

import (
	"encoding/json"
	"io"
	"os"
)

type args struct {
	Val string
}

type reply struct {
	Val   string
	Error error
}

func main() {
	jsWr := json.NewEncoder(os.Stdout)
	jsRd := json.NewDecoder(os.Stdin)
	in := &args{}
	err := jsRd.Decode(in)
	out := &reply{}
	if err != nil {
		out.Error = err
	} else {
		out.Val = in.Val
	}
	if err := jsWr.Encode(out); err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}
}
