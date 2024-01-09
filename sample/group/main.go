package main

import (
	"errors"
	"github.com/desertbit/grumble"
	"strings"
	"time"
)

var App = grumble.New(&grumble.Config{
	Name:        "foo",
	Description: "An awesome foo bar",
	Flags: func(f *grumble.Flags) {
		f.String("d", "directory", "DEFAULT", "set an alternative root directory path")
		f.Bool("v", "verbose", false, "enable verbose mode")
	},
})

func init() {
	App.AddCommand(&grumble.Command{
		Name:    "daemon",
		Help:    "run the daemon",
		Aliases: []string{"run"},
		Flags: func(f *grumble.Flags) {
			f.Duration("t", "timeout", time.Second, "timeout duration")
		},
		Args: func(a *grumble.Args) {
			a.Bool("production", "whether to start the daemon in production or development mode")
			a.Int("opt-level", "the optimization mode", grumble.Default(3))
			a.StringList("services", "additional services that should be started", grumble.Default([]string{"test", "te11"}))
		},
		Run: func(c *grumble.Context) error {
			c.App.Println("timeout:", c.Flags.Duration("timeout"))
			c.App.Println("directory:", c.Flags.String("directory"))
			c.App.Println("verbose:", c.Flags.Bool("verbose"))
			c.App.Println("production:", c.Args.Bool("production"))
			c.App.Println("opt-level:", c.Args.Int("opt-level"))
			c.App.Println("services:", strings.Join(c.Args.StringList("services"), ","))
			return nil
		},
	})

	adminCommand := &grumble.Command{
		Name:     "admin",
		Help:     "admin tools",
		LongHelp: "super administration tools",
	}
	App.AddCommand(adminCommand)

	adminCommand.AddCommand(&grumble.Command{
		Name: "root",
		Help: "root the machine",
		Run: func(c *grumble.Context) error {
			c.App.Println(c.Flags.String("directory"))
			return errors.New("failed")
		},
	})

	group := grumble.NewGroup("test group")
	group.AddCommand(&grumble.Command{
		Name: "group1",
		Help: "group1",
		Run: func(c *grumble.Context) error {
			c.App.Println(c.Flags.String("directory"))
			return errors.New("failed")
		},
	})
	g2cmd := &grumble.Command{
		Name: "group2",
		Help: "group2",
		Run: func(c *grumble.Context) error {
			c.App.Println(c.Flags.String("directory"))
			return errors.New("failed")
		},
	}
	g2cmd.AddCommand(&grumble.Command{
		Name: "subgroup2",
		Help: "subgroup2",
		Run: func(c *grumble.Context) error {
			c.App.Println(c.Flags.String("directory"))
			return errors.New("failed")
		},
	})
	group.AddCommand(g2cmd)
	App.AddGroup(group)

}

func main() {
	grumble.Main(App)

}
