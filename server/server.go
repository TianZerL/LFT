package server

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/TianZerL/LFT/headinfo"
)

//Server for receiving
type Server struct {
	IPAddr string
	Name   string
}

//set bufSize
const bufSize int = 4096

//sayHelloBack will send a response to the scanner
func sayHelloBack(addr, name string) {
	//Listing scanner
	lAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	l, err := net.ListenUDP("udp", lAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()
	buf := make([]byte, 64)
	for {
		//Waiting for scanner send a Hello
		n, rAddr, err := l.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		data := string(buf[:n])
		if data == "Hello" {
			//Send back server name
			_, err := l.WriteToUDP([]byte(name), rAddr)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func receiveHandler(conn net.Conn, dist string) {
	defer conn.Close()
	h := headinfo.NewHeadInfo()
	h.Waiting(conn)
	//Creat empty dir
	if h.IsDir == true {
		os.MkdirAll(dist+h.Name, os.ModePerm)
		return
	}
	//Creat all necessary dirs
	os.MkdirAll(dist+filepath.Dir(h.Name), os.ModePerm)
	path := setFileName(dist + h.Name)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		log.Println(err)
		return
	}
	buf := make([]byte, bufSize)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			log.Printf(` Path: "%s"  | Size: %.1f kb`, path, float64(h.Size)/1024.0)
			break
		}
		f.Write(buf[:n])
	}
}

//checkIP will check if it is legal of a IP
func checkIP(ip, port string) bool {
	if net.ParseIP(ip) != nil {
		if i, err := strconv.ParseInt(port, 10, 64); i > 2000 && i <= 65535 && err == nil {
			return true
		}
	}
	return false
}

//NewServer will Creat a server
func NewServer(name, ip, port string) (*Server, error) {
	if checkIP(ip, port) == false {
		return nil, errors.New("illegal IP addres")
	}
	return &Server{IPAddr: ip + ":" + port, Name: name}, nil
}

//Waiting for receiving
func (s *Server) Waiting(dist string) {
	dist += "/"
	go sayHelloBack(s.IPAddr, s.Name)
	l, err := net.Listen("tcp", s.IPAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()
	log.Println("Server start listening " + s.IPAddr)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go receiveHandler(conn, dist)
	}
}
