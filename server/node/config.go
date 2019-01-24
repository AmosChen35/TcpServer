package node

import(
    "fmt"
    "os"
    "strings"
    "path/filepath"

    "github.com/AmosChen35/TcpServer/server/rpc"
)

type Config struct {
    Name string `toml:"-"`
    NodeVersion string `toml:",omitempty"`
    TCPHost string `toml:",omitempty"`
    TCPPort int `toml:",omitempty"`
    TCPTimeouts rpc.TCPTimeouts
}

func (c *Config) name() string {
    if c.Name == "" {
        progname := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
        if progname == "" {
            panic("empty executable name, set Config.Name")
        }
        return progname
    }
    return c.Name
}

// HTTPEndpoint resolves an HTTP endpoint based on the configured host interface
// and port parameters.
func (c *Config) TCPEndpoint() string {
    if c.TCPHost == "" {
        return ""
    }
    return fmt.Sprintf("%s:%d", c.TCPHost, c.TCPPort)
}
