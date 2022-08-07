package main

import (
	"bufio"
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
}

var ErrCompleted = errors.New("completed")

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	reader := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	return &Client{
		address:   address,
		timeout:   timeout,
		reader:    reader,
		writer:    writer,
		in:        in,
		completed: false,
		b:         make([]byte, 1024),
	}
}

type Client struct {
	address   string
	timeout   time.Duration
	conn      net.Conn
	_conn     *bufio.Reader
	reader    *bufio.Scanner
	writer    *bufio.Writer
	in        io.ReadCloser
	completed bool
	b         []byte
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	c._conn = bufio.NewReader(c.conn)
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Send() (err error) {
	defer func() {
		if err != nil {
			c.completed = true
		}
	}()
	if c.completed {
		return ErrCompleted
	}
	if c.reader.Scan() {
		text := c.reader.Text()
		_, err = c.conn.Write([]byte(text + "\n"))
	} else {
		err = ErrCompleted
	}
	return
}

func (c *Client) Receive() (err error) {
	defer func() {
		if err != nil {
			c.completed = true
		}
	}()
	if c.completed {
		return ErrCompleted
	}
	n, err := c.conn.Read(c.b)
	if errors.Is(err, io.EOF) {
		log.Println("...EOF")
		return ErrCompleted
	}
	if err != nil {
		return err
	}
	_, err = c.writer.Write(c.b[:n])
	_ = c.writer.Flush()
	return
}
