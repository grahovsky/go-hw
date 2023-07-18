package main

import (
	"context"
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrNoConnection     = errors.New("NO CONNECTION")
	ErrClosedConnection = errors.New("CONNECTION IS CLOSED")
	ErrWrite            = errors.New("FAILED TO WRITE")
	ErrRead             = errors.New("FAILED TO READ")
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

	var err error
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if _, err = io.Copy(tc.conn, tc.in); err != nil {
				break OUTER
			}
		}
	}

	return err
}

func (tc *telnetClient) Receive(ctx context.Context) error {
	if tc.conn == nil {
		return ErrNoConnection
	}

	var err error

OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if _, errCopy := io.Copy(tc.out, tc.conn); errCopy != nil {
				if errors.Is(errCopy, io.EOF) {
					err = ErrClosedConnection
				} else {
					err = errCopy
				}
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
