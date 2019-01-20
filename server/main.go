package main

import(
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/AmosChen35/TcpServer/server/node"
    "github.com/AmosChen35/TcpServer/server/core"
    "github.com/AmosChen35/TcpServer/server/params"
)

func startNode(n *node.Node) error {
    if err := n.Start(); err != nil {
        return fmt.Errorf("Error starting protocol node: %v", err)
    }

    go func() {
        sigc := make(chan os.Signal, 1)
        signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
        defer signal.Stop(sigc)
        <-sigc
        fmt.Println("Got interrupt, shutting down...")
        go n.Stop()
        for i := 10; i > 0; i-- {
            <-sigc
            if i > 1 {
                fmt.Println("Already shutting down, interrupt more to panic.", "times", i-1)
            }
        }
    }()

    return nil
}

func main() {
    fmt.Println("Server Start")

    config := &node.Config{
        Name:    "Server1",
        NodeVersion: params.ChainVersion,
    }

    myNode, err := node.New(config)
    if err != nil {
        panic(err)
    }

    if err := myNode.Register(func(ctx *node.ServiceContext) (node.Service, error) {
        return core.New(ctx, &core.Config{
            Name:       "CoreService",
            TCPPort:    8080,
        })
    }); err != nil {
        panic(err)
    }

    if err := startNode(myNode); err != nil {
        panic(err)
    }

    //Main observer
    go func() {
        for {
            fmt.Println("OK")
            time.Sleep(time.Duration(1)*time.Second)
        }
    }()

    var core *core.Core
    if err := myNode.Service(&core); err != nil {
        panic(err)
    }

    myNode.Wait()
}
