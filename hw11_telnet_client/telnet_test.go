package main

import (
	"bytes"
	"context"
	"io"
	"net"
	"os/signal"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			ctx, _ := signal.NotifyContext(context.Background())

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect(ctx))
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)

			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("basic with new lines", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		msgSentByTelnet := "I\nam\nTELNET\nclient!!!\n"
		msgSentByClient := "hello\nworld\n"

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			ctx, _ := signal.NotifyContext(context.Background())

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect(ctx))
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString(msgSentByTelnet)
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, msgSentByClient, out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, msgSentByTelnet, string(request)[:n])

			n, err = conn.Write([]byte(msgSentByClient))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("error connection", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("2s")
		require.NoError(t, err)

		ctx, _ := signal.NotifyContext(context.Background())
		client := NewTelnetClient(":4000", timeout, io.NopCloser(in), out)
		require.Error(t, client.Connect(ctx))
	})
}
