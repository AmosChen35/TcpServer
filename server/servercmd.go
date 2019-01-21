package main

import(
    "os"
    "fmt"
    "runtime"

    "github.com/AmosChen35/TcpServer/server/utils"
    "github.com/AmosChen35/TcpServer/server/params"
    "github.com/AmosChen35/TcpServer/server/console"
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
