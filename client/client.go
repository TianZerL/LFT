package client

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/TianZerL/LFT/headinfo"
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

func (c *Client) sendFile(src string) {
	conn, err := net.Dial("tcp", c.IPAddr)
	if err != nil {
		log.Fatalln(err)
	}
	infoChan = make(chan string)
	log.Println("Start")
	mode := &headinfo.ModeInfo{
		Mode:     headinfo.ModeFile,
		BasePath: filepath.Dir(src)}
	go sendHandler(conn, src, mode)
	log.Println(<-infoChan)
	log.Println("All finished, exiting...")
}

func (c *Client) sendDir(src string) {
	files, dirs := getDirInfo(src)
	infoChan = make(chan string, 20)
	//Start send
	log.Println("Start")
	//establish directory structure in the server
	for _, dir := range dirs {
		conn, err := net.Dial("tcp", c.IPAddr)
		if err != nil {
			log.Fatalln(err)
		}
		mode := &headinfo.ModeInfo{
			Mode:     headinfo.ModeDir,
			BasePath: filepath.Dir(src)}
		establishHandler(conn, dir, mode)
	}
	//Waiting for establishing
	for range dirs {
		<-infoChan
	}
	log.Println("Directory structure established successfully")
	//Send files
	for _, file := range files {
		conn, err := net.Dial("tcp", c.IPAddr)
		if err != nil {
			log.Fatalln(err)
		}
		mode := &headinfo.ModeInfo{
			Mode:     headinfo.ModeDir,
			BasePath: filepath.Dir(src)}
		go sendHandler(conn, file, mode)
	}
	for range files {
		log.Println(<-infoChan)
	}
	log.Println("All finished, exiting...")
}

//Handler for establishing dirs in server
func establishHandler(conn net.Conn, dir string, mode *headinfo.ModeInfo) {
	//Error handling
	var err error = nil
	defer func() {
		if err != nil && err != io.EOF {
			infoChan <- fmt.Sprintf(`"%s": Fail to establish directory structure ! error: %v`, dir, err)
		}
	}()
	defer conn.Close()
	//Get ready to go
	dirInfo := headinfo.NewHeadInfo()
	ds, err := os.Stat(dir)
	if err != nil {
		return
	}
	dirInfo.Init(ds, dir, mode)
	err = dirInfo.Send(conn)
	if err != nil {
		return
	}
	//finished
	infoChan <- "f"
}

//Handler for sending file
func sendHandler(conn net.Conn, src string, mode *headinfo.ModeInfo) {
	//Error handling
	var err error = nil
	defer func() {
		if err != nil && err != io.EOF {
			infoChan <- fmt.Sprintf(`"%s": File send failed! error: %v`, src, err)
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
	h.Init(fs, src, mode)
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
