package service

import (
	"app/internal/app/core/port"
	"encoding/hex"
	"math/rand"
)

type Service struct {
	Tg       port.Itg
	Repo     port.IRepo
	RepoUser port.IRepoUser
	RepoAuth port.IRepoAuth
}

func New(
	tg port.Itg,
	repo port.IRepo,
	repoUser port.IRepoUser,
	repoAuth port.IRepoAuth,
) *Service {
	return &Service{
		Tg:       tg,
		Repo:     repo,
		RepoUser: repoUser,
		RepoAuth: repoAuth,
	}
}




func (s *Service) Token(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}