package scan

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Scanner defined the basic information to scan server
type Scanner struct {
	Ports []string
	IPs   []string
}

var lock sync.RWMutex
var wg sync.WaitGroup

//checkIP will check if it is legal of a IP
func checkIP(ip string) bool {
	if net.ParseIP(ip) != nil && ip != "0.0.0.0" {
		return true
	}
	return false
}

//checkIP will check if it is legal of a port
func checkPort(port string) bool {
	if i, err := strconv.ParseInt(port, 10, 64); i > 2000 && i <= 65535 && err == nil {
		return true
	}
	return false
}

//NewScanner will Creat a scanner
func NewScanner(ip, port string) (*Scanner, error) {
	ips := strings.Split(ip, ",")
	ports := strings.Split(port, ",")
	for _, ip := range ips {
		if checkIP(ip) == false {
			return nil, errors.New("illegal IP addres")
		}
	}
	for _, port := range ports {
		if checkPort(port) == false {
			return nil, errors.New("illegal port")
		}
	}
	return &Scanner{Ports: ports, IPs: ips}, nil
}

func sayHelloHandler(conn net.Conn) {
	wg.Add(1)
	defer wg.Done()
	defer conn.Close()
	Addr := conn.RemoteAddr()
	_, err := conn.Write([]byte{'H', 'e', 'l', 'l', 'o'})
	if err != nil {
		log.Println(err)
		return
	}
	callBack := make([]byte, 64)
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Println(err)
		return
	}
	n, err := conn.Read(callBack)
	conn.RemoteAddr()
	if err != nil {
		log.Println(err)
		return
	}
	lock.Lock()
	fmt.Println("Server name:   " + string(callBack[:n]) + "      IP:   " + Addr.String())
	lock.Unlock()
}

//Scan the specified IP segments and ports
func (s *Scanner) Scan() {
	//IP
	log.Println("Start scan servers")
	for _, ip := range s.IPs {
		ipInfo := strings.SplitAfter(ip, ".")
		//port
		for _, port := range s.Ports {
			//Like 192.168.1.0 to 192.168.1.255
			for i := int64(0); i < 256; i++ {
				addr := ipInfo[0] + ipInfo[1] + ipInfo[2] + strconv.FormatInt(i, 10) + ":" + port
				conn, err := net.Dial("udp", addr)
				if err != nil {
					log.Fatalln(err)
				}
				go sayHelloHandler(conn)
			}
		}
	}
	wg.Wait()
	log.Println("All finished, exiting...")
}
