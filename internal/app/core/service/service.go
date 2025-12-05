package service

import (
	"app/internal/app/core/port"
)

type Service struct {
	salt     []byte
	Tg       port.Itg
	MapAuth  port.IMapAuth
	Repo     port.IRepo
	RepoUser port.IRepoUser
	RepoAuth port.IRepoAuth
}

func New(
	salt string,
	tg port.Itg,
	mapAuth port.IMapAuth,
	repo port.IRepo,
	repoUser port.IRepoUser,
	repoAuth port.IRepoAuth,
) *Service {
	return &Service{
		salt:     []byte(salt),
		Tg:       tg,
		MapAuth:  mapAuth,
		Repo:     repo,
		RepoUser: repoUser,
		RepoAuth: repoAuth,
	}
}
