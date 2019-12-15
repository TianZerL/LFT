package server

import (
	"LFT/headinfo"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

//Server for receiving
type Server struct {
	IPAddr string
}

//set bufSize
const bufSize int = 4096

//TODO
func receiveHandler(conn net.Conn, dist string) {
	defer conn.Close()
	h := headinfo.NewHeadInfo()
	h.Waiting(conn)
	os.MkdirAll(dist, os.ModePerm)
	if h.IsDir == true {
		//TODO
		return
	}
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
func NewServer(ip, port string) (*Server, error) {
	if checkIP(ip, port) == false {
		return nil, errors.New("illegal IP addres")
	}
	return &Server{IPAddr: ip + ":" + port}, nil
}

//Waiting for receiving
func (s *Server) Waiting(dist string) {
	dist += "/"
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
		}
		go receiveHandler(conn, dist)
	}
}
