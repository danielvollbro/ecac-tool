package hclmodels

type RootRaw struct {
	Hosts   []*Host   `hcl:"host,block"`
	Plugins []*Plugin `hcl:"plugin,block"`
	Tasks   []*Task   `hcl:"task,block"`
}
