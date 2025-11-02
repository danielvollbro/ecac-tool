package hclmodels

type Root struct {
	Hosts   map[string]*Host   `hcl:"host,block"`
	Plugins map[string]*Plugin `hcl:"plugin,block"`
	Tasks   map[string]*Task   `hcl:"task,block"`
}
