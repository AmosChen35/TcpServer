package utils

import (
    "encoding"
    "flag"
    "fmt"
    "os"
    "os/user"
    "path"
    "strings"

    "gopkg.in/urfave/cli.v1"
)

// Custom type which is registered in the flags library which cli uses for
// argument parsing. This allows us to expand Value to an absolute path when
// the argument is parsed
type DirectoryString struct {
    Value string
}

func (self *DirectoryString) String() string {
    return self.Value
}

func (self *DirectoryString) Set(value string) error {
    self.Value = expandPath(value)
    return nil
}

// Custom cli.Flag type which expand the received string to an absolute path.
// e.g. ~/.ethereum -> /home/username/.ethereum
type DirectoryFlag struct {
    Name  string
    Value DirectoryString
    Usage string
}

func (self DirectoryFlag) String() string {
    fmtString := "%s %v\t%v"
    if len(self.Value.Value) > 0 {
        fmtString = "%s \"%v\"\t%v"
    }
    return fmt.Sprintf(fmtString, prefixedNames(self.Name), self.Value.Value, self.Usage)
}

func eachName(longName string, fn func(string)) {
    parts := strings.Split(longName, ",")
    for _, name := range parts {
        name = strings.Trim(name, " ")
        fn(name)
    }
}

// called by cli library, grabs variable from environment (if in env)
// and adds variable to flag set for parsing.
func (self DirectoryFlag) Apply(set *flag.FlagSet) {
    eachName(self.Name, func(name string) {
        set.Var(&self.Value, self.Name, self.Usage)
    })
}

type TextMarshaler interface {
    encoding.TextMarshaler
    encoding.TextUnmarshaler
}

// textMarshalerVal turns a TextMarshaler into a flag.Value
type textMarshalerVal struct {
    v TextMarshaler
}

func (v textMarshalerVal) String() string {
    if v.v == nil {
        return ""
    }
    text, _ := v.v.MarshalText()
    return string(text)
}

func (v textMarshalerVal) Set(s string) error {
    return v.v.UnmarshalText([]byte(s))
}

// TextMarshalerFlag wraps a TextMarshaler value.
type TextMarshalerFlag struct {
    Name  string
    Value TextMarshaler
    Usage string
}

func (f TextMarshalerFlag) GetName() string {
    return f.Name
}

func (f TextMarshalerFlag) String() string {
    return fmt.Sprintf("%s \"%v\"\t%v", prefixedNames(f.Name), f.Value, f.Usage)
}

func (f TextMarshalerFlag) Apply(set *flag.FlagSet) {
    eachName(f.Name, func(name string) {
        set.Var(textMarshalerVal{f.Value}, f.Name, f.Usage)
    })
}

// GlobalTextMarshaler returns the value of a TextMarshalerFlag from the global flag set.
func GlobalTextMarshaler(ctx *cli.Context, name string) TextMarshaler {
    val := ctx.GlobalGeneric(name)
    if val == nil {
        return nil
    }
    return val.(textMarshalerVal).v
}

func prefixFor(name string) (prefix string) {
    if len(name) == 1 {
        prefix = "-"
    } else {
        prefix = "--"
    }

    return
}

func prefixedNames(fullName string) (prefixed string) {
    parts := strings.Split(fullName, ",")
    for i, name := range parts {
        name = strings.Trim(name, " ")
        prefixed += prefixFor(name) + name
        if i < len(parts)-1 {
            prefixed += ", "
        }
    }
    return
}

func (self DirectoryFlag) GetName() string {
    return self.Name
}

func (self *DirectoryFlag) Set(value string) {
    self.Value.Value = value
}

// Expands a file path
// 1. replace tilde with users home dir
// 2. expands embedded environment variables
// 3. cleans the path, e.g. /a/b/../c -> /a/c
// Note, it has limitations, e.g. ~someuser/tmp will not be expanded
func expandPath(p string) string {
    if strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
        if home := homeDir(); home != "" {
            p = home + p[1:]
        }
    }
    return path.Clean(os.ExpandEnv(p))
}

func homeDir() string {
    if home := os.Getenv("HOME"); home != "" {
        return home
    }
    if usr, err := user.Current(); err == nil {
        return usr.HomeDir
    }
    return ""
}
