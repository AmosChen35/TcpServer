package node

import(
    "os"
    "strings"
    "path/filepath"
)

type Config struct {
    Name string `toml:"-"`
    NodeVersion string `toml:",omitempty"`
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
