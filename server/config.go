package main

import(
    "github.com/AmosChen35/TcpServer/server/node"
    "github.com/AmosChen35/TcpServer/server/core"
    "github.com/AmosChen35/TcpServer/server/params"
    "gopkg.in/urfave/cli.v1"
)

func makeNode(ctx *cli.Context) *node.Node{
    config := &node.Config{
        Name:    "Server1",
        NodeVersion: params.Version,
        TCPHost: "127.0.0.1",
        TCPPort: 8080,
    }

    myNode, err := node.New(config)
    if err != nil {
        panic(err)
    }

    if err := myNode.Register(func(ctx *node.ServiceContext) (node.Service, error) {
        return core.New(ctx, &core.Config{
            Name:       "CoreService",
            TCPPort:    8081,
        })
    }); err != nil {
        panic(err)
    }

    return myNode
}
