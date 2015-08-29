package net

import (
	"bufio"
	"net"

	"github.com/golang/glog"
)

type RunConn struct {
	InComing chan string
	OutGoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func NewRunConn(conn net.Conn) *RunConn {
	rc := &RunConn{
		InComing: make(chan string),
		OutGoing: make(chan string),
		reader:   bufio.NewReader(conn),
		writer:   bufio.NewWriter(conn),
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
		if _, err := rc.writer.WriteString(data); err != nil {
			glog.Errorf("write >> WriteString: %s", err.Error())
			close(rc.OutGoing)
			return
		}
		if err := rc.writer.Flush(); err != nil {
			glog.Errorf("write >> Flush: %s", err.Error())
			close(rc.OutGoing)
			return
		}
	}
}

func (rc *RunConn) read() {
	var line string
	var err error
	for {
		if line, err = rc.reader.ReadString('\n'); err != nil {
			glog.Errorf("read >> ReadString: %s", err.Error())
			close(rc.InComing)
			return
		}
		rc.InComing <- line
	}
}
