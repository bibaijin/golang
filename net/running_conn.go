package net

import (
	"bufio"
	"net"

	"github.com/golang/glog"
)

type RunningConn struct {
	InComing chan string
	OutGoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func NewRunningConn(conn net.Conn) *RunningConn {
	rc := &RunningConn{
		InComing: make(chan string),
		OutGoing: make(chan string),
		reader:   bufio.NewReader(conn),
		writer:   bufio.NewWriter(conn),
	}

	rc.run()

	return rc
}

func (rc *RunningConn) run() {
	go rc.write()
	go rc.read()
}

func (rc *RunningConn) write() {
	for data := range rc.OutGoing {
		if _, err := rc.writer.WriteString(data); err != nil {
			glog.Error(err.Error())
			return
		}
		if err := rc.writer.Flush(); err != nil {
			glog.Error(err.Error())
			return
		}
	}
}

func (rc *RunningConn) read() {
	var line string
	var err error
	for {
		if line, err = rc.reader.ReadString('\n'); err != nil {
			glog.Error(err.Error())
			break
		}
		rc.InComing <- line
	}
}
