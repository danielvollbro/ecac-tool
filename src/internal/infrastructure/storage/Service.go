package storage

import "os"

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) FileExists(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
