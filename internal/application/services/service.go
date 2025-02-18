package services

import (
	repos "github.com/Jereyji/auth-service.git/internal/domain/interface_repository"
	"github.com/Jereyji/auth-service.git/internal/pkg/configs"
)

type Service struct {
	repository repos.RepositoryI
	trm        repos.TransactionManagerI
	config     *configs.AuthConfig
}

func NewService(rep repos.RepositoryI, config *configs.AuthConfig, trm repos.TransactionManagerI) *Service {
	return &Service{
		repository: rep,
		trm:        trm,
		config:     config,
	}
}
