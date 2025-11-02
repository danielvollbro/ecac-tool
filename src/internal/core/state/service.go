package state

import (
	"ecac/internal/models"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Service struct {
	FilePath string
}

func NewService(filePath string) *Service {
	return &Service{
		FilePath: filePath,
	}
}

func (s *Service) SaveState(state models.State) {
	state.RunAt = time.Now()

	data, _ := json.MarshalIndent(state, "", "  ")
	path, _ := os.Getwd()
	file := filepath.Join(path, s.FilePath)

	if err := os.WriteFile(file, data, 0644); err != nil {
		log.Fatalf("Kunde inte spara state: %v", err)
	}
}

func (s *Service) LoadState() models.State {
	path, _ := os.Getwd()
	file := filepath.Join(path, s.FilePath)

	data, err := os.ReadFile(file)
	if err != nil {
		return models.State{}
	}

	var state models.State
	if err := json.Unmarshal(data, &state); err != nil {
		return models.State{}
	}
	return state
}
