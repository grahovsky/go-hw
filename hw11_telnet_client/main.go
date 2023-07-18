package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Println("Usage: go-telnet --timeout=10s host port")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	tc := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := tc.Connect(ctx); err != nil {
		fmt.Printf("Connection error: %v. %v \n\r", address, err)
		os.Exit(1)
	}
	defer tc.Close()

	go func() {
		if err := tc.Receive(ctx); err != nil {
			log.Fatal(err)
		}
		stop()
	}()

	go func() {
		if err := tc.Send(ctx); err != nil {
			log.Fatal(err)
		}
		stop()
	}()

	<-ctx.Done()
	stop()
}
