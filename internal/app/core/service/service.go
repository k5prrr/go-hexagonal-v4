package service

import (
	"app/internal/app/core/domain"
	"app/internal/app/core/port"
	"context"
	"encoding/hex"
	"fmt"
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
func (s *Service) Code() string  {
	codeInt := rand.Intn(899999) + 100000

	return fmt.Sprintf("%d", codeInt)
}


func (s *Service) UserFull(ctx context.Context, id int64) (*domain.UserFull, error) {
/*	user, err := s.RepoUser.Get(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	auth, err := s.RepoAuth.GetByInt(ctx, "user_id", user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, auth, nil

 */
	return nil, nil
}
func (s *Service) UserAuthByPhone(ctx context.Context, phone string) (*domain.User, *domain.Auth, error) {
	user, err := s.RepoUser.GetBy(ctx, "phone", phone)
	if err != nil {
		return nil, nil, err
	}

	auth, err := s.RepoAuth.GetByInt(ctx, "user_id", user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, auth, nil
}