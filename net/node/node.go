package node

import (
	"crypto/rand"
	"encoding/hex"
	"net"
)

type Node struct {
	Host string
	Port int
	Stat NodeStat
	Id   [8]byte
}

func NewNode(host string, port int) Node {
	node := Node{
		Host: host,
		Port: port,
	}
	rand.Read(node.Id[:])
	return node
}

func NewNodeWithID(host string, port int, id [8]byte) Node {
	return Node{
		Host: host,
		Port: port,
		Id:   id,
	}
}

func NewNodeWithAddress(address *net.TCPAddr) Node {
	return Node{
		Host: address.IP.String(),
		Port: address.Port,
	}
}

func (n Node) GetHexId() string {
	return hex.EncodeToString(n.Id[:])
}

func (n Node) GetAddress() net.TCPAddr {
	return net.TCPAddr{
		IP:   net.ParseIP(n.Host),
		Port: n.Port,
	}
}

func (n Node) Equals(o Node) bool {
	return n.Host == o.Host && n.Port == o.Port
}
