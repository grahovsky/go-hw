package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	ErrNoConnection     = errors.New("NO CONNECTION")
	ErrClosedConnection = errors.New("CONNECTION IS CLOSED")
)

type TelnetClient interface {
	Connect(context.Context) error
	io.Closer
	Send(context.Context) error
	Receive(context.Context) error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (tc *telnetClient) Connect(ctx context.Context) error {
	dialer := net.Dialer{Timeout: tc.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", tc.address)
	if err != nil {
		return err
	}
	tc.conn = conn
	return nil
}

func (tc *telnetClient) Send(ctx context.Context) error {
	if tc.conn == nil {
		return ErrNoConnection
	}

	scanner := bufio.NewScanner(tc.in)

	var err error

OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if scanner.Scan() {
				str := scanner.Text()
				if str == "quit" {
					break OUTER
				}
				tc.conn.Write([]byte(fmt.Sprintf("%s\n", str)))
			}
		}
	}
	return err
}

func (tc *telnetClient) Receive(ctx context.Context) error {
	if tc.conn == nil {
		return ErrNoConnection
	}

	scanner := bufio.NewScanner(tc.conn)
	var err error

OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if scanner.Scan() {
				str := scanner.Text()
				tc.out.Write([]byte(fmt.Sprintf("%s\n", str)))
			} else {
				err = ErrClosedConnection
				break OUTER
			}
		}
	}

	return err
}

func (tc *telnetClient) Close() error {
	if tc.conn == nil {
		return ErrNoConnection
	}
	return tc.conn.Close()
}
