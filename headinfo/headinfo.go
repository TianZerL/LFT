package headinfo

import (
	"bytes"
	"net"
	"os"
	"strconv"
	"time"
)

//HeadInfo defined a file's head infomation
type HeadInfo struct {
	Name  string
	Size  int64
	IsDir bool
}

//NewHeadInfo return a new headinfo type
func NewHeadInfo() *HeadInfo {
	return &HeadInfo{}
}

//Init initiation a headinfo type
func (h *HeadInfo) Init(f os.FileInfo) {
	h.Name = f.Name()
	h.Size = f.Size()
	h.IsDir = f.IsDir()
}

//Send headinfo to server
func (h *HeadInfo) Send(conn net.Conn) error {
	data := make([]byte, 0, 1024)
	data = strconv.AppendBool(append(strconv.AppendInt(append(append(data, []byte(h.Name)...), 0x01), h.Size, 10), 0x01), h.IsDir)
	_, err := conn.Write(data)
	if err != nil {
		return err
	}
	re := make([]byte, 6)
	err = conn.SetReadDeadline(time.Now().Add(20 * time.Second))
	if err != nil {
		return err
	}
	_, err = conn.Read(re)
	if err != nil {
		return err
	}
	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		return err
	}
	return nil
}

//Waiting for headinfo
func (h *HeadInfo) Waiting(conn net.Conn) error {
	data := make([]byte, 1024*16)
	n, err := conn.Read(data)
	if err != nil {
		return err
	}
	info := bytes.Split(data[:n], []byte{0x01})
	h.Name = string(info[0])
	h.Size, err = strconv.ParseInt(string(info[1]), 10, 64)
	if err != nil {
		return err
	}
	h.IsDir, err = strconv.ParseBool(string(info[2]))
	if err != nil {
		return err
	}
	conn.Write([]byte{'S', 't', 'a', 'r', 't'})
	return nil
}
