package main

import (
	"log"

	"github.com/TianZerL/LFT/client"
	"github.com/TianZerL/LFT/options"
	"github.com/TianZerL/LFT/scan"
	"github.com/TianZerL/LFT/server"
)

func main() {
	opt := options.NewOptions()
	if opt.Help == true {
		opt.Usage()
	} else if opt.Server == true {
		s, err := server.NewServer(opt.Name, opt.IP, opt.Port)
		if err != nil {
			log.Fatalln(err)
		}
		s.Waiting(opt.Dist)
	} else if opt.Scan == true {
		scanner, err := scan.NewScanner(opt.IP, opt.Port)
		if err != nil {
			log.Fatalln(err)
		}
		scanner.Scan()
	} else {
		c, err := client.NewClient(opt.IP, opt.Port)
		if err != nil {
			log.Fatalln(err)
		}
		c.Send(opt.Dist)
	}
}
