package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
	Start(ctx context.Context, stop context.CancelFunc)
}

var ErrNotConnected = errors.New("not connected")

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *client) Close() error {
	var err error
	if c.conn != nil {
		err = c.conn.Close()
		c.conn = nil
	}
	return err
}

func (c *client) Send() error {
	if c.conn == nil {
		return ErrNotConnected
	}
	_, err := io.Copy(c.conn, c.in)
	return err
}

func (c *client) Receive() error {
	if c.conn == nil {
		return ErrNotConnected
	}
	_, err := io.Copy(c.out, c.conn)
	return err
}

func (c *client) Start(ctx context.Context, stop context.CancelFunc) {
	c.startReader(ctx, stop)
	c.startWriter(stop)
}

func (c *client) startReader(ctx context.Context, stop context.CancelFunc) {
	go func() {
		scan := bufio.NewScanner(c.in)
		for scan.Scan() {
			if _, err := c.conn.Write([]byte(scan.Text() + "\n")); err != nil {
				log.Println(err)
				stop()
				return
			}
		}
		select {
		case <-ctx.Done():
			return
		default:
			if _, err := c.conn.Write([]byte("Bye-bye\n")); err != nil {
				log.Println(err)
			}
			stop()
		}
	}()
}

func (c *client) startWriter(stop context.CancelFunc) {
	go func() {
		for {
			if err := c.Receive(); err != nil {
				log.Println(err)
				stop()
				return
			}
		}
	}()
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
