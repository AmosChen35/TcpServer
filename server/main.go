package main

import(
    "fmt"
    "os"
    "sort"
    "os/signal"
    "syscall"
    "time"

    "github.com/AmosChen35/TcpServer/server/node"
    "github.com/AmosChen35/TcpServer/server/core"
    "github.com/AmosChen35/TcpServer/server/utils"
    "gopkg.in/urfave/cli.v1"
)

var (
    // Git SHA1 commit hash of the release (set via linker flags)
    gitCommit = ""
    // The app that holds all commands and flags.
    app = utils.NewApp(gitCommit, "the go-ethereum command line interface")
    // flags that configure the node
    nodeFlags = []cli.Flag{
        utils.TestFlag,
        utils.DataDirFlag,
    }
)

func init() {
    app.Action = Server
    app.HideVersion = true // we have a command to print the version
    app.Copyright = "Copyright 2019 ..."
    app.Commands = []cli.Command{
        consoleCommand,
        versionCommand,
        licenseCommand,
    }
    sort.Sort(cli.CommandsByName(app.Commands))

    app.Flags = append(app.Flags, nodeFlags...)

    app.Before = func(ctx *cli.Context) error {
        return nil
    }

    app.After = func(ctx *cli.Context) error {
        return nil
    }
}

func startNode(ctx *cli.Context, n *node.Node) error {
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

    //Main observer
    go func() {
        for {
            time.Sleep(time.Duration(1)*time.Second)
        }
    }()

    var core *core.Core
    if err := n.Service(&core); err != nil {
        panic(err)
    }

    return nil
}

func Server(ctx *cli.Context) error {
    if args := ctx.Args(); len(args) > 0 {
        return fmt.Errorf("invalid command: %q", args[0])
    }

    fmt.Println("Server Start")
    node := makeNode(ctx)
    if err := startNode(ctx, node); err != nil {
        panic(err)
    }
    node.Wait()
    return nil
}

func main() {
    if err := app.Run(os.Args); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
