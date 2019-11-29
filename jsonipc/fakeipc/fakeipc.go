package main

import (
	"encoding/json"
	"io"
	"os"
)

type args struct {
	Arg string
}

type reply struct {
	Result string
	Error  error
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
		out.Result = in.Arg
	}
	if err := jsWr.Encode(out); err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}
}
