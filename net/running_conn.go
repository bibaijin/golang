package net

import (
	"io"
	"net"

	"github.com/golang/glog"
)

type RunConn struct {
	InComing  chan []byte
	OutGoing  chan []byte
	conn      net.Conn
	isRunning chan bool
}

func NewRunConn(conn net.Conn) *RunConn {
	rc := &RunConn{
		InComing:  make(chan []byte, 1),
		OutGoing:  make(chan []byte, 1),
		conn:      conn,
		isRunning: make(chan bool, 1),
	}

	rc.run()

	return rc
}

func (rc *RunConn) run() {
	go rc.write()
	go rc.read()
}

func (rc *RunConn) write() {
	defer glog.V(2).Info("write >> done")

	for data := range rc.OutGoing {
		if _, err := rc.conn.Write(data); err != nil {
			// glog.Errorf("write >> %s", err.Error())
			return
		}
	}
}

func (rc *RunConn) read() {
	defer glog.V(2).Info("read >> done")

	tmp := make([]byte, 256)
	for {
		n, err := rc.conn.Read(tmp)
		if err != nil || n == 0 {
			if err != io.EOF {
				// glog.Error(err.Error())
			}
			rc.stop()
			return
		}

		data := make([]byte, n)
		copy(data, tmp[0:n])

		select {
		case rc.InComing <- data:
		case <-rc.isRunning:
			return
		}
	}
}

func (rc *RunConn) stop() {
	close(rc.InComing)
}

func (rc *RunConn) Close() {
	rc.isRunning <- false
	rc.conn.Close()
}
