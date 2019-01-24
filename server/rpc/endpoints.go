package rpc

import(
    "fmt"
    "net"
)

func StartTCPEndpoint(endpoint string, apis []API, timeouts TCPTimeouts) (net.Listener, *Server, error) {
    handler := NewServer()
    for _, api := range apis {
        if api.Public {
            if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
                return nil, nil, err
            }
            fmt.Println("TCP registered", "service", api.Service, "namespace", api.Namespace)
        }
    }
    var(
        listener net.Listener
        err error
    )
    if listener, err = net.Listen("tcp", endpoint); err != nil {
        return nil, nil, err
    }
    go NewTCPServer(listener, handler, timeouts)
    return listener, handler, err
}
