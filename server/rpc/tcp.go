package rpc

import(
    "fmt"
    "net"
    "time"
    "context"
)

type TCPTimeouts struct {
    ReadTimeout time.Duration
    WriteTimeout time.Duration
    IdleTimeout time.Duration
}

var DefaultTCPTimeouts = TCPTimeouts{
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
    IdleTimeout:  120 * time.Second,
}

func handleConnection(conn net.Conn) (chan error){
    defer conn.Close()
    notify := make(chan error)


    return notify
}

func NewTCPServer(listener net.Listener, handler *Server, timouts TCPTimeouts) {
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println(err.Error())
        }

        handler.codecsMu.Lock()
        c := handler.AddConnection(conn)
        handler.codecsMu.Unlock()

        fmt.Println("RemoteAddr:", conn.RemoteAddr())
        go handler.ServeCodec(NewJSONCodec(conn, conn.RemoteAddr(), c), OptionMethodInvocation|OptionSubscriptions)
    }
}

func DialTCP(endpoint string) (*Client, error) {
    initctx := context.Background()
    c, err := newClient(initctx, func(context.Context) (net.Conn, error) {
        fmt.Println(endpoint)
        conn, err := net.Dial("tcp", endpoint)
        if err != nil {
            return nil, err
        }
        return conn, nil
    })
    return c, err
}

