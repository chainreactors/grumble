package grumble

import (
	"fmt"
)

func NewGroup(name string) *Group {
	return &Group{
		Name:     name,
		commands: &Commands{},
		enabled:  true,
	}
}

type Group struct {
	Name     string
	commands *Commands
	enabled  bool
}

func (g *Group) validate() error {
	if g.Name == "" {
		return fmt.Errorf("empty group name")
	}
	return nil
}

func (g *Group) Commands() *Commands {
	return g.commands
}

func (g *Group) AddCommand(cmd *Command) {
	err := cmd.validate()
	if err != nil {
		panic(err)
	}
	cmd.parent = nil
	cmd.registerFlagsAndArgs(true)
	g.commands.Add(cmd)
}

func (g *Group) Enable() {
	g.enabled = true
}

func (g *Group) Disable() {
	g.enabled = false
}

type Groups map[string]*Group

func (g Groups) Find(name string) *Group {
	return g[name]
}

func (g Groups) Parse(args []string, parentFlagMap FlagMap, skipFlagMaps bool) (cmds []*Command, flagsMap FlagMap, rest []string, err error) {
	for _, i := range g {
		if !i.enabled {
			continue
		}
		cmds, flagsMap, args, err = i.commands.parse(args, parentFlagMap, skipFlagMaps)
		if err != nil {
			return
		} else if len(cmds) > 0 {
			return
		}
	}
	return nil, nil, nil, fmt.Errorf("not found any command")
}

func (g Groups) FindCommand(args []string) (cmd *Command, rest []string, err error) {
	var cmds []*Command
	cmds, _, rest, err = g.Parse(args, nil, true)
	if err != nil {
		return
	}

	if len(cmds) > 0 {
		cmd = cmds[len(cmds)-1]
	}

	return
}

func (g Groups) Commands() *Commands {
	var cmds Commands
	for _, i := range g {
		if !i.enabled {
			continue
		}
		cmds.list = append(cmds.list, i.commands.list...)
	}
	return &cmds
}

func (g Groups) Enable(name string) bool {
	if i := g[name]; i == nil {
		return false
	} else {
		i.Enable()
		return true
	}
}

func (g Groups) Disable(name string) bool {
	if i := g[name]; i == nil {
		return false
	} else {
		i.Disable()
		return true
	}
}
