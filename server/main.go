package main

import(
    "fmt"

    "github.com/AmosChen35/TcpServer/server/node"
    "github.com/AmosChen35/TcpServer/server/core"
    "github.com/AmosChen35/TcpServer/server/params"
)

func main() {
    fmt.Println("Server Start")

    config := &node.Config{
        Name:    "Server1",
        NodeVersion: params.ChainVersion,
    }

    myNode, err := node.New(config)
    if err != nil { panic(err)}

    if err := myNode.Register(func(ctx *node.ServiceContext) (node.Service, error) {
        return core.New(ctx, &core.Config{
            Name:       "CoreService",
            TCPPort:    8080,
        })
    }); err != nil {
        panic(err)
    }

    myNode.Start()

    var core *core.Core
    if err := myNode.Service(&core); err != nil {
        panic(err)
    }
}
