package hcl

import (
	"ecac/internal/infrastructure/storage"
	hclmodels "ecac/internal/models/hcl-models"
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Service struct {
	Storage storage.Service
}

func NewService(storage storage.Service) *Service {
	return &Service{
		Storage: storage,
	}
}

func (s *Service) Load(path string) (*hclmodels.Root, error) {
	exists := s.Storage.FileExists(path)
	if !exists {
		return nil, fmt.Errorf("%s: file doesn't exists", path)
	}

	var rootRaw hclmodels.RootRaw
	err := hclsimple.DecodeFile(path, nil, &rootRaw)
	if err != nil {
		return nil, err
	}

	config := &hclmodels.Root{
		Hosts:   make(map[string]*hclmodels.Host),
		Plugins: make(map[string]*hclmodels.Plugin),
		Tasks:   make(map[string]*hclmodels.Task),
	}

	for _, h := range rootRaw.Hosts {
		config.Hosts[h.Name] = h
	}

	for _, p := range rootRaw.Plugins {
		config.Plugins[p.Name] = p
	}

	for _, t := range rootRaw.Tasks {
		config.Tasks[t.Name] = t
	}

	return config, nil
}
