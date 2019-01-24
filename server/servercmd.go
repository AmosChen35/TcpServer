package main

import(
    "os"
    "fmt"
    "runtime"
    "strings"

    "github.com/AmosChen35/TcpServer/server/utils"
    "github.com/AmosChen35/TcpServer/server/params"
    "github.com/AmosChen35/TcpServer/server/console"
    "github.com/AmosChen35/TcpServer/server/rpc"
    "gopkg.in/urfave/cli.v1"
)

var(
    consoleCommand = cli.Command{
        Action:   utils.MigrateFlags(localConsole),
        Name:     "console",
        Usage:    "Start an interactive JavaScript environment",
        Flags:    nodeFlags,
        Category: "CONSOLE COMMANDS",
        Description: `
`,
    }
    attachCommand = cli.Command{
        Action:    utils.MigrateFlags(remoteConsole),
        Name:      "attach",
        Usage:     "Start an interactive JavaScript environment (connect to node)",
        ArgsUsage: "[endpoint]",
        Flags:     nodeFlags,
        Category:  "CONSOLE COMMANDS",
        Description: `
`,
    }
    versionCommand = cli.Command{
        Action:    utils.MigrateFlags(version),
        Name:      "version",
        Usage:     "Print version numbers",
        ArgsUsage: " ",
        Category:  "MISCELLANEOUS COMMANDS",
        Description: `
The output of this command is supposed to be machine-readable.
`,
    }
    licenseCommand = cli.Command{
        Action:    utils.MigrateFlags(license),
        Name:      "license",
        Usage:     "Display license information",
        ArgsUsage: " ",
        Category:  "MISCELLANEOUS COMMANDS",
    }
)

func localConsole(ctx *cli.Context) error {
    // Create and start the node based on the CLI flags
    node := makeNode(ctx)
    if err := startNode(ctx, node); err != nil {
        panic(err)
    }
    defer node.Stop()

    client, err := node.Attach()
    if err != nil {
        fmt.Printf("Failed to attach to the inproc geth: %v", err)
    }

    config := console.Config{
        DataDir: utils.MakeDataDir(ctx),
        //DocRoot: ctx.GlobalString(utils.JSpathFlag.Name),
        Client:  client,
        //Preload: utils.MakeConsolePreloads(ctx),
    }

    console, err := console.New(config)
    if err != nil {
        fmt.Printf("Failed to start the JavaScript console: %v", err)
    }
    defer console.Stop(false)

    // If only a short execution was requested, evaluate and return
    if script := ctx.GlobalString(utils.ExecFlag.Name); script != "" {
        console.Evaluate(script)
        return nil
    }
    // Otherwise print the welcome screen and enter interactive mode
    console.Welcome()
    console.Interactive()

    return nil
}

func remoteConsole(ctx *cli.Context) error {
    // Attach to a remotely running geth instance and start the JavaScript console
    endpoint := ctx.Args().First()
    if endpoint == "" {
        panic("the remote endpoint string missing.")
    }
    client, err := dialRPC(endpoint)
    if err != nil {
        return fmt.Errorf("Unable to attach to remote server: %v", err)
    }

    config := console.Config{
        DataDir: utils.MakeDataDir(ctx),
        //DocRoot: ctx.GlobalString(utils.JSpathFlag.Name),
        Client:  client,
        //Preload: utils.MakeConsolePreloads(ctx),
    }

    console, err := console.New(config)
    if err != nil {
        fmt.Printf("Failed to start the JavaScript console: %v", err)
    }
    defer console.Stop(false)

    if script := ctx.GlobalString(utils.ExecFlag.Name); script != "" {
        console.Evaluate(script)
        return nil
    }

    // Otherwise print the welcome screen and enter interactive mode
    console.Welcome()
    console.Interactive()

    return nil
}

func dialRPC(endpoint string) (*rpc.Client, error) {
    if endpoint == "" {
        panic("the remote endpoint string missing.")
    } else if strings.HasPrefix(endpoint, "tcp:") {
        endpoint = endpoint[4:]
    }
    return rpc.Dial(endpoint)
}

func version(ctx *cli.Context) error {
    fmt.Println("Version:", params.VersionWithMeta)
    if gitCommit != "" {
        fmt.Println("Git Commit:", gitCommit)
    }
    fmt.Println("Architecture:", runtime.GOARCH)
    fmt.Println("Go Version:", runtime.Version())
    fmt.Println("Operating System:", runtime.GOOS)
    fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
    fmt.Printf("GOROOT=%s\n", runtime.GOROOT())
    return nil
}

func license(_ *cli.Context) error {
    fmt.Println(`You should have received a copy of the GNU General Public License. If not, see <http://www.gnu.org/licenses/>.`)
    return nil
}
