package main

import (
	"flag"
	"goroutine-example/constant"
	"log"
)

type InputParameterParser struct {
	poolCap   int
	host      string
	network   string
	startPort int
	endPort   int
}

func (p *InputParameterParser) PoolCap() int {
	return p.poolCap
}

func (p *InputParameterParser) SetPoolCap(poolCap int) {
	p.poolCap = poolCap
}

func (p *InputParameterParser) Host() string {
	return p.host
}

func (p *InputParameterParser) SetHost(host string) {
	p.host = host
}

func (p *InputParameterParser) Network() string {
	return p.network
}

func (p *InputParameterParser) SetNetwork(network string) {
	p.network = network
}

func (p *InputParameterParser) StartPort() int {
	return p.startPort
}

func (p *InputParameterParser) SetStartPort(startPort int) {
	p.startPort = startPort
}

func (p *InputParameterParser) EndPort() int {
	return p.endPort
}

func (p *InputParameterParser) SetEndPort(endPort int) {
	p.endPort = endPort
}

func GetInputParameterParser() *InputParameterParser {
	return &InputParameterParser{
		poolCap:   constant.PoolCap,
		host:      constant.Host,
		network:   constant.Network,
		startPort: constant.StartPort,
		endPort:   constant.EndPort,
	}
}

func InputParse() *InputParameterParser {

	parser := GetInputParameterParser()

	poolCap := flag.Int("poolCap", constant.PoolCap, "number of goroutine")
	host := flag.String("host", constant.Host, "host: eg 127.0.0.1、127.0.0.1/24")
	network := flag.String("network", constant.Network, "network type: eg tcp、icmp")
	startPort := flag.Int("startPort", constant.StartPort, "host: eg 127.0.0.1、127.0.0.1/24")
	endPort := flag.Int("endPort", constant.EndPort, "host: eg 127.0.0.1、127.0.0.1/24")
	flag.Parse()

	for i := 0; i != flag.NArg(); i++ {
		log.Printf("arg[%d]=%s\n", i, flag.Arg(i))
	}

	parser.poolCap = *poolCap
	parser.host = *host
	parser.startPort = *startPort
	parser.endPort = *endPort
	parser.network = *network

	return parser
}
