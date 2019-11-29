package main

import (
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"time"
)

type Args struct {
	Val string
}

type Reply struct {
	Result string
	Error  error
}

type stdioRWC struct {
	r io.ReadCloser
	w io.WriteCloser
}

func (rwc stdioRWC) Read(p []byte) (n int, err error) {
	println("read")
	return rwc.r.Read(p)
}

func (rwc stdioRWC) Write(p []byte) (n int, err error) {
	println("write")
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
	println("crap")
	reply.Result = args.Val
	return nil
}

func main() {
	rwc := stdioRWC{r: os.Stdin, w: os.Stdout}

	srv := rpc.NewServer()
	srv.RegisterName("Echo", &handler{})
	srv.ServeCodec(jsonrpc.NewServerCodec(rwc))
	time.Sleep(1 * time.Second)
}
