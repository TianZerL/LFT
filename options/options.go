package options

import (
	"flag"
)

//Opt defined options LFT can accept
type Opt struct {
	IP     string
	Port   string
	Dist   string
	Server bool
	Scan   bool
	Help   bool
}

var options Opt

func init() {
	flag.BoolVar(&options.Server, "w", false, "Start a server")
	flag.BoolVar(&options.Scan, "scan", false, "Scan Lan to find servers(TODO)")
	flag.BoolVar(&options.Help, "h", false, "Display help information")
	flag.BoolVar(&options.Help, "?", false, "Display help information")
	flag.StringVar(&options.IP, "ip", "0.0.0.0", "Server IP address")
	flag.StringVar(&options.Port, "port", "6981", "Server Port")
	flag.StringVar(&options.Dist, "d", "./receive/", "Source or destination")
	flag.Parse()
}

//NewOptions return a options object, which is been initialized.
func NewOptions() *Opt {
	return &options
}

//Usage print help infomation
func (o *Opt) Usage() {
	flag.PrintDefaults()
}
