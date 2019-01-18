package core

import(
)

var DefaultConfig = Config{
    TCPPort:    1234,
}

type Config struct{
    Name    string      `toml:"-"`
    TCPPort int         `toml:",omitempty"`
}
