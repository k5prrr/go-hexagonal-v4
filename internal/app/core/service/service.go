package service

import (
	"app/internal/app/core/port"
)

type Service struct {
	Repo     port.IRepo
	RepoUser port.IRepoUser
	RepoAuth port.IRepoAuth
}

func New(
	repo port.IRepo,
	repoUser port.IRepoUser,
	repoAuth port.IRepoAuth,
	) *Service {
	return &Service{
		Repo: repo,
		RepoUser: repoUser,
		RepoAuth: repoAuth,
	}
}
