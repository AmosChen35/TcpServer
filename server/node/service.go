package node

import(
    "fmt"
    "reflect"
    "bytes"
    "github.com/BurntSushi/toml"
)

type ServiceContext struct{
    config         *Config
    services       map[reflect.Type]Service // Index of the already constructed services
}

type Service interface {
    Start() error
    Stop() error
}

func (ctx *ServiceContext) DoSomethingWithContext(tcpPort int) error {
    buf := new(bytes.Buffer)
    if err := toml.NewEncoder(buf).Encode(ctx.config); err != nil {
        return err
    }

    fmt.Printf("Core Service Port = %d\n", tcpPort)
    fmt.Println(buf)

    return nil
}
