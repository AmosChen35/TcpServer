package rpc

import(
    "time"
)

type TCPTimeouts struct {
    ReadTimeout time.Duration
    WriteTimeout time.Duration
    IdleTimeout time.Duration
}

var DefaultTCPTimeouts = TCPTimeouts{
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
    IdleTimeout:  120 * time.Second,
}

func NewTCPServer() {

}
