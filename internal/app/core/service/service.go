package service

import (
	"app/internal/app/core/port"
)

type Service struct {
	Repo     port.IRepo
	RepoUser port.IRepoUser
}

func New(
	repo port.IRepo,
	repoUser port.IRepoUser,
	) *Service {
	return &Service{
		Repo: repo,
		RepoUser: repoUser,
	}
}
