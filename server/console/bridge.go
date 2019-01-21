package console

import (
    "fmt"
    "io"

    "github.com/AmosChen35/TcpServer/server/rpc"
)

// bridge is a collection of JavaScript utility methods to bride the .js runtime
// environment and the Go RPC connection backing the remote method calls.
type bridge struct {
    client   *rpc.Client  // RPC client to execute Ethereum requests through
    prompter UserPrompter // Input prompter to allow interactive user feedback
    printer  io.Writer    // Output writer to serialize any display strings to
}

// newBridge creates a new JavaScript wrapper around an RPC client.
func newBridge(client *rpc.Client, prompter UserPrompter, printer io.Writer) *bridge {
    return &bridge{
        client:   client,
        prompter: prompter,
        printer:  printer,
    }
}

func (b *bridge) HelloBridge() string {
    fmt.Println("OK")
    return "HelloBridge"
}
