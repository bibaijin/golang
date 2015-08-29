package net

import (
	"net"
	"time"
)

type TimeoutConn struct {
	timeout time.Duration
	conn    net.Conn
}

func NewTimeoutConn(timeout time.Duration, conn net.Conn) *TimeoutConn {
	tc := &TimeoutConn{
		timeout: timeout,
		conn:    conn,
	}
	return tc
}

func (tc *TimeoutConn) Read(b []byte) (n int, err error) {
	tc.conn.SetReadDeadline(time.Now().Add(tc.timeout))
	n, err = tc.conn.Read(b)
	return n, err
}

func (tc *TimeoutConn) Write(b []byte) (n int, err error) {
	tc.conn.SetWriteDeadline(time.Now().Add(tc.timeout))
	n, err = tc.conn.Write(b)
	return n, err
}

func (tc *TimeoutConn) Close() error {
	return tc.conn.Close()
}

func (tc *TimeoutConn) LocalAddr() net.Addr {
	return tc.conn.LocalAddr()
}

func (tc *TimeoutConn) RemoteAddr() net.Addr {
	return tc.conn.RemoteAddr()
}
