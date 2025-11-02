package hclmodels

import "github.com/hashicorp/hcl/v2"

type Task struct {
	Name   string   `hcl:"name,label"`
	Plugin string   `hcl:"plugin"`
	Body   hcl.Body `hcl:",remain"`
	Config hcl.Body
}
