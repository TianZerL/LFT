package main

import (
	"LFT/client"
	"LFT/options"
	"LFT/server"
	"log"
)

func main() {
	opt := options.NewOptions()
	if opt.Help == true {
		opt.Usage()
		return
	}
	if opt.Server == true {
		s, err := server.NewServer(opt.IP, opt.Port)
		if err != nil {
			log.Fatalln(err)
		}
		s.Waiting(opt.Dist)
	} else {
		c, err := client.NewClient(opt.IP, opt.Port)
		if err != nil {
			log.Fatalln(err)
		}
		c.Send(opt.Dist)
	}
}
