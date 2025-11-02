package sshagent

import "golang.org/x/crypto/ssh"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Connect() error {
	var hostKey ssh.PublicKey

	config := &ssh.ClientConfig{
		User: "daniel",
		Auth: []ssh.AuthMethod{
			ssh.Password("123qweASD"),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	client, err := ssh.Dial("tcp", "192.168.50.111", config)
	if err != nil {
		return err
	}
	defer client.Close()

	return nil
}
