package core

import (
    "fmt"

    "github.com/AmosChen35/TcpServer/server/node"
    "github.com/AmosChen35/TcpServer/server/rpc"
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
    fmt.Printf("[%s] start\n", core.config.Name)
    return nil
}

func (core *Core) Stop() error {
    fmt.Printf("[%s] stop\n", core.config.Name)
    return nil
}

func (core *Core) HelloCore() map[string]string {
    fmt.Printf("[%s] api HelloCore\n", core.config.Name)

    hello := make(map[string]string)
    hello["WELCOME RPC"] = "1.0"
    hello["TEST"] = core.config.Name

    return hello
}

func (core *Core) APIs() []rpc.API {
    return []rpc.API{
        {
            Namespace: "core",
            Version:   "3.0",
            Service:   core,
            Public:    true,
        },
    }
}
