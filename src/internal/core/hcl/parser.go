package hcl

import (
	"ecac/internal/infrastructure/storage"
	hclmodels "ecac/internal/models/hcl-models"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type Service struct {
	Storage storage.Service
}

func NewService(storage storage.Service) *Service {
	return &Service{
		Storage: storage,
	}
}

func (s *Service) ParseConfig(path string) (*hclmodels.RootRaw, error) {
	exists := s.Storage.FileExists(path)
	if !exists {
		return nil, fmt.Errorf("%s: file doesn't exists", path)
	}

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, diags
	}

	var rootRaw hclmodels.RootRaw
	ctx := &hcl.EvalContext{}
	diags = gohcl.DecodeBody(file.Body, ctx, &rootRaw)
	if diags.HasErrors() {
		return nil, diags
	}

	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "config"},
		},
	}

	for _, t := range rootRaw.Tasks {
		content, _ := t.Body.Content(schema)
		if len(content.Blocks) == 0 {
			return nil, fmt.Errorf("%s: task %q missing config block", path, t.Name)
		}
		t.Config = content.Blocks[0].Body
	}

	return &rootRaw, nil
}
