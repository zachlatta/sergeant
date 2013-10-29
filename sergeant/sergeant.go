package sergeant

import (
  "flag"
  "fmt"
  "io"
  "os"
  "strings"
  "sync"
  "text/template"
  "unicode"
  "unicode/utf8"
)

// A command is an implementation of a command-line command issued by the user,
// like 'git status'
type Command struct {
  // Run runs the command.
  // The args are the arguments after the command name.
  Run func(cmd *Command, args []string)

  // UsageLine is the one-line usage message.
  // The first word in the line is taken to be the command name.
  UsageLine string

  // Short is the short description shown in the 'help' output.
  Short string

  // Long is the long message shown in the 'help <this-command>' output.
  Long string

  // Flag is a set of flags specific to this command
  Flag flag.FlagSet

  // CustomFlags indicates that the command will do its own flag parsing.
  CustomFlags bool
}

// Name returns the command's name: the first word in the command's usage line.
func (c *Command) Name() string {
  name := c.UsageLine
  i := strings.Index(name, " ")
  if i >= 0 {
    name = name[:i]
  }
  return name
}


// Usage prints the command's usage to stderr and exists the program with exit
// code 2
func (c *Command) Usage() {
  fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
  fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
  os.Exit(2)
}

// Runnable repotrs whether the command can be run; otherwise it is a
// documentatino pseudo-command.
func (c *Command) Runnable() bool {
  return c.Run != nil
}

var exitStatus = 0
var exitMu sync.Mutex

func setExitStatus(n int) {
  exitMu.Lock()
  if exitStatus < n {
    exitStatus = n
  }
  exitMu.Unlock()
}


// The list of commands in the application. The user will need populate this in
// their application.
var Commands = []*Command{}

// The name of the application. The user will need to set this.
var ApplicationName string

// A brief description of the application. The user will need to set this.
var ApplicationDescription string

func Init() {
  flag.Usage = usage
  flag.Parse()

  args := flag.Args()
  if len(args) < 1 {
    usage()
  }

  if args[0] == "help" {
    help(args[1:])
    return
  }

  for _, cmd := range Commands {
    if cmd.Name() == args[0] && cmd.Run != nil {
      cmd.Flag.Usage = func() { cmd.Usage() }
      if cmd.CustomFlags {
        args = args[1:]
      } else {
        cmd.Flag.Parse(args[1:])
        args = cmd.Flag.Args()
      }
      cmd.Run(cmd, args)
      exit()
      return
    }
  }

  errMsg := ApplicationName + ": unknown subcommand %q\nRun '" + ApplicationName + " help` for usage.\n"

  fmt.Fprintf(os.Stderr, errMsg, args[0])
}

var usageTemplate = `{{.Description}}

Usage:

  {{.Name}} command [arguments]

The commands are:
{{range .Data}}{{if .Runnable}}
  {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "{{.Name}} help [command]" for more information about a command.

`

var helpTemplate = `{{if .Data.Runnable}}usage: {{.Name}} {{.Data.UsageLine}}

{{end}}{{.Data.Long | trim}}

`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
  t := template.New("top")
  t.Funcs(template.FuncMap{
    "trim": strings.TrimSpace,
    "capitalize": capitalize,
  })
  template.Must(t.Parse(text))

  templateData := struct {
    Data interface{}
    Name string
    Description string
  } {
    data,
    ApplicationName,
    ApplicationDescription,
  }
  if err := t.Execute(w, templateData); err != nil {
    panic(err)
  }
}

func capitalize(s string) string {
  if s == "" {
    return s
  }
  r, n := utf8.DecodeRuneInString(s)
  return string(unicode.ToTitle(r)) + s[n:]
}

func printUsage(w io.Writer) {
  tmpl(w, usageTemplate, Commands)
}

func usage() {
  printUsage(os.Stderr)
  os.Exit(2)
}

// help implements the 'help' command.
func help(args []string) {
  if len(args) == 0 {
    printUsage(os.Stdout)
    // not exit 2: succeeded at 'programname help'
    return
  }
  if len(args) != 1 {
    errMsg := "usage: " + ApplicationName + " help command\n\nToo many arguments given.\n"
    fmt.Fprintf(os.Stderr, errMsg)
    os.Exit(2) // failed at 'programname help'
  }

  arg := args[0]

  for _, cmd := range Commands {
    if cmd.Name() == arg {
      tmpl(os.Stdout, helpTemplate, cmd)
      // not exit 2: succeeded at 'programname help cmd'.
      return
    }
  }

  errMsg := "Unknown help topic %#q. Run '" + ApplicationName + " help'.\n"
  fmt.Fprintf(os.Stderr, errMsg, arg)
  os.Exit(2) // failed at 'programname help cmd'
}

var atexitFuncs []func()

func atexit(f func()) {
  atexitFuncs = append(atexitFuncs, f)
}

func exit() {
  for _, f := range atexitFuncs {
    f()
  }
  os.Exit(exitStatus)
}
