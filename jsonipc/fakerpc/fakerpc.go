package main

import (
	"fmt"
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

type JsonRPCError struct {
	Code    int
	Message string
	Data    interface{}
}

func (err JsonRPCError) Error() string {
	return err.Message
}

func (err JsonRPCError) MarshalJSON() ([]byte, error) {
	return []byte(`{"Code":1, "Message": "fake"}`), nil
}

type Args struct {
	Val string
}

type Reply struct {
	Val string
}

type stdioRWC struct {
	r io.ReadCloser
	w io.WriteCloser
}

func (rwc stdioRWC) Read(p []byte) (n int, err error) {
	return rwc.r.Read(p)
}

func (rwc stdioRWC) Write(p []byte) (n int, err error) {
	return rwc.w.Write(p)
}

func (rwc stdioRWC) Close() error {
	rErr := rwc.r.Close()
	wErr := rwc.w.Close()
	if rErr != nil {
		return rErr
	}
	return wErr
}

type handler struct {
}

func (h handler) Echo(args *Args, reply *Reply) error {
	reply.Val = args.Val
	return nil
}

func (h handler) Error(args *Args, reply *Reply) error {
	return fmt.Errorf("error with msg: %s", args.Val)
}

func (h handler) Error2(args *Args, reply *Reply) error {
	return JsonRPCError{Code: 1, Message: "error with msg", Data: args.Val}
}

func main() {
	rwc := stdioRWC{r: os.Stdin, w: os.Stdout}
	srv := rpc.NewServer()
	srv.RegisterName("Echo", &handler{})
	srv.ServeCodec(jsonrpc.NewServerCodec(rwc))
}
