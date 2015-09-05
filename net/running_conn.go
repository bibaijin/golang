package net

import (
	"io"
	"net"

	"github.com/golang/glog"
)

type RunConn struct {
	InComing chan []byte
	OutGoing chan []byte
	conn     net.Conn
}

func NewRunConn(conn net.Conn) *RunConn {
	rc := &RunConn{
		InComing: make(chan []byte),
		OutGoing: make(chan []byte),
		conn:     conn,
	}

	rc.run()

	return rc
}

func (rc *RunConn) run() {
	go rc.write()
	go rc.read()
}

func (rc *RunConn) write() {
	for data := range rc.OutGoing {
		if _, err := rc.conn.Write(data); err != nil {
			glog.Errorf("write >> WriteString: %s", err.Error())
			close(rc.OutGoing)
			return
		}
	}
}

func (rc *RunConn) read() {
	tmp := make([]byte, 256)
	for {
		if _, err := rc.conn.Read(tmp); err != nil {
			// glog.Errorf("read >> ReadString: %s", err.Error())
			if err != io.EOF {
				glog.Error(err.Error())
			}
			close(rc.InComing)
			return
		}
		rc.InComing <- tmp
	}
}
