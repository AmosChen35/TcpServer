package utils

import(
    "os"
    "fmt"
    "path/filepath"

    "github.com/AmosChen35/TcpServer/server/params"
    "github.com/AmosChen35/TcpServer/server/node"
    "gopkg.in/urfave/cli.v1"
)

var (
    CommandHelpTemplate = `{{.cmd.Name}}{{if .cmd.Subcommands}} command{{end}}{{if .cmd.Flags}} [command options]{{end}} [arguments...]
{{if .cmd.Description}}{{.cmd.Description}}
{{end}}{{if .cmd.Subcommands}}
SUBCOMMANDS:
    {{range .cmd.Subcommands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
    {{end}}{{end}}{{if .categorizedFlags}}
{{range $idx, $categorized := .categorizedFlags}}{{$categorized.Name}} OPTIONS:
{{range $categorized.Flags}}{{"\t"}}{{.}}
{{end}}
{{end}}{{end}}`
)

func init() {
    cli.AppHelpTemplate = `{{.Name}} {{if .Flags}}[global options] {{end}}command{{if .Flags}} [command options]{{end}} [arguments...]

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

    cli.CommandHelpTemplate = CommandHelpTemplate
}

// NewApp creates an app with sane defaults.
func NewApp(gitCommit, usage string) *cli.App {
    app := cli.NewApp()
    app.Name = filepath.Base(os.Args[0])
    app.Author = ""
    app.Email = ""
    app.Version = params.VersionWithMeta
    if len(gitCommit) >= 8 {
        app.Version += "-" + gitCommit[:8]
    }
    app.Usage = usage
    return app
}

var (
	DataDirFlag = DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: DirectoryString{node.DefaultDataDir()},
	}
    TestFlag = cli.StringFlag{
        Name: "Test",
        Usage: `Usage Test ("start", "stop")`,
        Value: "full",
    }
    ExecFlag = cli.StringFlag{
        Name:  "exec",
        Usage: "Execute JavaScript statement",
    }
)

func MigrateFlags(action func(ctx *cli.Context) error) func(*cli.Context) error {
    return func(ctx *cli.Context) error {
        for _, name := range ctx.FlagNames() {
            if ctx.IsSet(name) {
                ctx.GlobalSet(name, ctx.String(name))
            }
        }
        return action(ctx)
    }
}

// MakeDataDir retrieves the currently requested data directory, terminating
// if none (or the empty string) is specified. If the node is starting a testnet,
// the a subdirectory of the specified datadir will be used.
func MakeDataDir(ctx *cli.Context) string {
	if path := ctx.GlobalString(DataDirFlag.Name); path != "" {
		return path
	}
	fmt.Println("Cannot determine default data directory, please set manually (--datadir)")
	return ""
}
