package headinfo

import (
	"bytes"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//SendMode defined the mode to send
type SendMode uint8

const (
	//ModeFile to send a file
	ModeFile SendMode = 1
	//ModeDir to send a Dir
	ModeDir SendMode = 2
)

//ModeInfo defined the file's send mode, including baed path
type ModeInfo struct {
	Mode     SendMode
	BasePath string
}

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
func (h *HeadInfo) Init(f os.FileInfo, src string, mode *ModeInfo) {
	if mode.Mode == ModeDir {
		relPath, err := filepath.Rel(mode.BasePath, src)
		if err != nil {
			log.Fatalln(err)
		}
		h.Name = relPath
	} else {
		h.Name = f.Name()
	}
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
