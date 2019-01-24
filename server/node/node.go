package node

import (
    "net"
    "fmt"
    "sync"
    "reflect"
    "strings"
    "errors"

    "github.com/AmosChen35/TcpServer/server/rpc"
)

type Node struct{
    config *Config

    serviceFuncs []ServiceConstructor     // Service constructors (in dependency order)
    services     map[reflect.Type]Service // Currently running services

    rpcAPIs       []rpc.API   // List of APIs currently provided by the node
    inprocHandler *rpc.Server // In-process RPC request handler to process the API requests

    tcpEndpoint string  //TCP endpoint (IP + PORT)
    tcpListener net.Listener  //TCP listener socket to server API request
    tcpHandler *rpc.Server // TCP request handler to process the API request

    running bool   // node running flag
    stop    chan struct{} // Channel to wait for termination notifications

    lock sync.RWMutex
}

type ServiceConstructor func(ctx *ServiceContext) (Service, error)

// Before new node do some config check.
func New(conf *Config) (*Node, error) {

    if strings.ContainsAny(conf.Name, `/\`) {
        return nil, errors.New(`Config name must not content '\' or '/'`)
    }

    return &Node {
        config: conf,
        tcpEndpoint:      conf.TCPEndpoint(),
    }, nil
}

func (n *Node) Register(constructor ServiceConstructor) error {
    n.lock.Lock()
    defer n.lock.Unlock()

    n.serviceFuncs = append(n.serviceFuncs, constructor)
    return nil
}

// Service retrieves a currently running service registered of a specific type.
func (n *Node) Service(service interface{}) error {
    n.lock.RLock()
    defer n.lock.RUnlock()

    // Otherwise try to find the service to return
    element := reflect.ValueOf(service).Elem()
    if running, ok := n.services[element.Type()]; ok {
        element.Set(reflect.ValueOf(running))
        return nil
    }
    return ErrServiceUnknown
}

func (n *Node) Start() error {
    n.lock.Lock()
    defer n.lock.Unlock()

    n.running = true

    services := make(map[reflect.Type]Service)
    for _, constructor := range n.serviceFuncs {
        ctx := &ServiceContext{
            config:         n.config,
            services:       make(map[reflect.Type]Service),
        }
        for kind, s := range services { // copy needed for threaded access
            ctx.services[kind] = s
        }
        // Construct and save the service
        service, err := constructor(ctx)
        if err != nil {
            return err
        }
        kind := reflect.TypeOf(service)
        if _, exists := services[kind]; exists {
            return &DuplicateServiceError{Kind: kind}
        }
        services[kind] = service
    }

    // Start each of the services
    started := []reflect.Type{}
    for kind, service := range services {
        // Start the next service, stopping all previous upon failure
        if err := service.Start(); err != nil {
            for _, kind := range started {
                services[kind].Stop()
            }

            return err
        }
        // Mark the service started for potential cleanup
        started = append(started, kind)
    }

    // Lastly start the configured RPC interfaces
    if err := n.startRPC(services); err != nil {
        for _, service := range services {
            service.Stop()
        }
        return err
    }

    // Finish initializing the startup
    n.services = services
    n.stop = make(chan struct{})

    return nil
}

func (n *Node) Stop() error {
    n.lock.Lock()
    defer n.lock.Unlock()

    if n.running == false {
        return ErrNodeStopped
    }

    failure := &StopError {
        Services: make(map[reflect.Type]error),
    }
    for kind, service := range n.services {
        if err := service.Stop(); err != nil {
            failure.Services[kind] = err
        }
    }

    // unblock n.Wait
    n.services = nil
    n.running = false
    close(n.stop)

    return nil
}

func (n *Node) Wait() {
    n.lock.RLock()
    if n.running == false {
        n.lock.RUnlock()
        return
    }
    stop := n.stop
    n.lock.RUnlock()

    <-stop
}

func (n *Node) Restart() error {
    if err := n.Stop(); err != nil {
        return err
    }
    if err := n.Start(); err != nil {
        return err
    }
    return nil
}

func (n *Node) startTCP(endpoint string, apis []rpc.API, timeouts rpc.TCPTimeouts) error {
    if endpoint == "" {
        return nil
    }
    listener, handler, err := rpc.StartTCPEndpoint(endpoint, apis, timeouts)
    if err != nil {
        return err
    }
    fmt.Printf("TCP endpoint open %v\n", endpoint)

    n.tcpEndpoint = endpoint
    n.tcpListener = listener
    n.tcpHandler = handler

    return nil
}

func (n *Node) stopTCP() {
    if n.tcpListener != nil {
        n.tcpListener.Close()
        n.tcpListener = nil

        fmt.Printf("TCP endpoint closed %v\n", n.tcpEndpoint)
    }
    if n.tcpHandler != nil {
        n.tcpHandler.Stop()
        n.tcpHandler = nil
    }
}

// startRPC is a helper method to start all the various RPC endpoint during node
// startup. It's not meant to be called at any time afterwards as it makes certain
// assumptions about the state of the node.
func (n *Node) startRPC(services map[reflect.Type]Service) error {
    // Gather all the possible APIs to surface
    apis := n.apis()
    for _, service := range services {
        apis = append(apis, service.APIs()...)
    }
    // Start the various API endpoints, terminating all in case of errors
    if err := n.startInProc(apis); err != nil {
        return err
    }
    if err := n.startTCP(n.tcpEndpoint, apis, n.config.TCPTimeouts); err != nil {
        n.stopTCP()
        return err
    }
    // All API endpoints started successfully
    n.rpcAPIs = apis
    return nil
}

func (n *Node) startInProc(apis []rpc.API) error {
    // Register all the APIs exposed by the services
    handler := rpc.NewServer()
    for _, api := range apis {
        if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
            return err
        }
        fmt.Println("InProc registered", "namespace", api.Namespace)
    }
    n.inprocHandler = handler
    return nil
}

// Attach creates an RPC client attached to an in-process API handler.
func (n *Node) Attach() (*rpc.Client, error) {
    n.lock.RLock()
    defer n.lock.RUnlock()

    if n.running == false {
        return nil, ErrNodeStopped
    }
    return rpc.DialInProc(n.inprocHandler), nil
}

// retrive the current TCP endpoint used by the protocol stack
func (n *Node) TCPEndpoint() string {
    n.lock.RLock()
    defer n.lock.RUnlock()

    if n.tcpListener == nil {
        return n.tcpListener.Addr().String()
    }

    return n.tcpEndpoint
}

// apis returns the collection of RPC descriptors this node offers.
func (n *Node) apis() []rpc.API {
    return []rpc.API{
        {
            Namespace: "admin",
            Version:   "1.0",
            Service:   NewPrivateAdminAPI(n),
            Public:    true,
        },
    }
}
