package node

import (
    "sync"
    "reflect"
    "strings"
    "errors"
)

type Node struct{
    config *Config

    serviceFuncs []ServiceConstructor     // Service constructors (in dependency order)
    services     map[reflect.Type]Service // Currently running services

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

    // Finish initializing the startup
    n.services = services

    return nil
}
