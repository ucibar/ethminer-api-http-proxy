package main

import (
	"fmt"
	"net"
)

type ETHMinerConn struct {
	net.Conn
}

func ETHMinerDial(address string) (*ETHMinerConn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &ETHMinerConn{conn}, nil
}

func (eth *ETHMinerConn) Write(p []byte) (n int, err error) {
	return fmt.Fprintln(eth.Conn, string(p))
}
