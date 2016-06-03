// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package botgoram

import "github.com/Patrolavia/godot"

func (f *fsm) StateMap(name string) (dot string) {
	fix := func(n string) string {
		if n == "" {
			return name
		}
		return n
	}
	add := func(a, b string) string {
		if a != "" {
			b = a + `\n` + b
		}
		return b
	}
	g := godot.NewGraph(false, true, "StateMap", nil)
	nodes := make(map[string]godot.Node)
	nodes[name] = g.AddNode(name, map[string]string{"fillcolor": "#ccccff", "style": "filled"})
	node := func(n string) godot.Node {
		ret, ok := nodes[n]
		if !ok { // undocumented state, warning user by draw it red
			ret = g.AddNode(n, map[string]string{"bgcolor": "red"})
			nodes[n] = ret
		}
		return ret
	}

	for _, s := range f.sm {
		n := fix(s.Name())
		opt := make(map[string]string)
		label := n
		enter, leave := s.Actions()
		if enter != nil {
			label += `\n------------\n` + "enter action"
		}
		if leave != nil {
			label += `\n------------\n` + "leave action"
		}
		opt["label"] = label
		nodes[n] = g.AddNode(n, opt)
	}

	for _, s := range f.sm {
		n := fix(s.Name())
		for _, t := range s.Transitors() {
			n := n
			m := fix(t.State)
			opt := make(map[string]string)
			label := t.Desc
			switch {
			case t.IsHidden:
				opt["style"] = "dotted"
				n, m = m, n
			case t.IsFallback:
				label = add(label, "fallback")
			case t.Command != "" && t.Type == TextMsg:
				label = add(label, "Command: "+t.Command)
			default:
				label = add(label, t.Type)
			}
			opt["label"] = label

			g.AddEdge([]godot.Node{node(m), node(n)}, opt)
		}
	}

	return g.String()
}
