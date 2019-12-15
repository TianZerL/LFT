package client

import (
	"LFT/headinfo"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

//Client for send files
type Client struct {
	IPAddr string
}

//set bufSize
const bufSize int = 4096

//wait until all thread to exit
var infoChan chan string

//checkIP will check if it is legal of a IP
func checkIP(ip, port string) bool {
	if net.ParseIP(ip) != nil && ip != "0.0.0.0" {
		if i, err := strconv.ParseInt(port, 10, 64); i > 2000 && i <= 65535 && err == nil {
			return true
		}
	}
	return false
}

func checkDir(src string) (bool, error) {
	f, err := os.Stat(src)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}

func (c *Client) sendFile(src string) {
	conn, err := net.Dial("tcp", c.IPAddr)
	if err != nil {
		log.Fatalln(err)
		return
	}
	infoChan = make(chan string)
	log.Println("Start")
	go sendHandler(conn, src)
	log.Println(<-infoChan)
	log.Println("All finished, exiting...")
}

//TODO
func (c *Client) sendDir(src string) {
	//TODO
}

func sendHandler(conn net.Conn, src string) {
	//Error handling
	var err error = nil
	defer func() {
		if err != nil && err != io.EOF {
			infoChan <- fmt.Sprintf(`"%s": file send failed! error: %v`, src, err)
		}
	}()
	defer conn.Close()
	//Open file
	f, err := os.Open(src)
	if err != nil {
		return
	}
	defer f.Close()
	h := headinfo.NewHeadInfo()
	//fs is the fileinfo for f,which is needed to initliza headinfo
	fs, err := f.Stat()
	if err != nil {
		return
	}
	//Get file's headinfo
	h.Init(fs)
	//Send headinfo to server
	err = h.Send(conn)
	if err != nil {
		return
	}
	//creat a buffer
	buf := make([]byte, bufSize)
	//Sending
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err != io.EOF {
				return
			}
			break
		}
		conn.Write(buf[:n])
	}
	//Send a successful infomation
	infoChan <- fmt.Sprintf(`"%s": file send successful!`, src)
}

//NewClient will Creat a client
func NewClient(ip, port string) (*Client, error) {
	if checkIP(ip, port) == false {
		return nil, errors.New("illegal IP addres")
	}
	return &Client{IPAddr: ip + ":" + port}, nil
}

//Send will send src to server
func (c *Client) Send(src string) {
	src, err := filepath.Abs(filepath.Clean(src))
	if err != nil {
		log.Fatalln(err)
	}
	isDir, err := checkDir(src)
	if err != nil {
		log.Fatalln(err)
	}
	if isDir == true {
		c.sendDir(src)
	} else {
		c.sendFile(src)
	}
}
