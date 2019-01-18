package core

import (
    "fmt"

    "github.com/AmosChen35/TcpServer/server/node"
)

type Core struct {
    config *Config
}

func New(ctx *node.ServiceContext, config *Config) (*Core, error) {
    err := DoSomething(ctx, config)
    if err != nil {
        return nil, err
    }

    core := Core {
        config: config,
    }

    return &core, nil
}

func DoSomething(ctx *node.ServiceContext, config *Config) error {
    err := ctx.DoSomethingWithContext(config.TCPPort)
    if err != nil {
        return err
    }
    return nil
}

func (core *Core) Start() error {
    fmt.Printf("[%s] start", core.config.Name)
    return nil
}

func (core *Core) Stop() error {
    fmt.Println("[%s] stop", core.config.Name)
    return nil
}
