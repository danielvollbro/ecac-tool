package hclmodels

import "github.com/hashicorp/hcl/v2"

type Task struct {
	Name   string   `hcl:"name,label"`
	Plugin string   `hcl:"plugin"`
	Config hcl.Body `hcl:"config,block"`
	Range  hcl.Range
}
