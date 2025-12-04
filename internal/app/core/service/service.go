package service

import (
	"app/internal/app/core/domain"
	"app/internal/app/core/port"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Service struct {
	salt     string
	Tg       port.Itg
	Repo     port.IRepo
	RepoUser port.IRepoUser
	RepoAuth port.IRepoAuth
}

func New(
	salt string,
	tg port.Itg,
	repo port.IRepo,
	repoUser port.IRepoUser,
	repoAuth port.IRepoAuth,
) *Service {
	return &Service{
		salt:     salt,
		Tg:       tg,
		Repo:     repo,
		RepoUser: repoUser,
		RepoAuth: repoAuth,
	}
}

func (s *Service) Secret(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
func (s *Service) Code() string {
	codeInt := rand.Intn(899999) + 100000

	return fmt.Sprintf("%d", codeInt)
}

func (s *Service) Sha256(secret, data string) string {
	data = fmt.Sprintf("%s%s", data, s.salt)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sig := h.Sum(nil)

	return base64.RawURLEncoding.EncodeToString(sig)
}
func (s *Service) Token(userID int64, secret string) string {
	payload := fmt.Sprintf("%d|%d", userID, time.Now().Add(30*24*time.Hour).Unix())
	sha := s.Sha256(secret, payload)
	tokenRaw := fmt.Sprintf("%s|%s", payload, sha)

	return base64.RawURLEncoding.EncodeToString([]byte(tokenRaw))
}
func (s *Service) CurrentUser(ctx context.Context, token string) (*domain.UserFull, error) {
	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token encoding: %w", err)
	}

	parts := bytes.Split(data, []byte{'|'})
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	userIDRaw, expRaw, sigB64 := parts[0], parts[1], parts[2]

	userID, err := strconv.ParseInt(string(userIDRaw), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	exp, err := strconv.ParseInt(string(expRaw), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid expiration time: %w", err)
	}

	if time.Now().Unix() > exp {
		return nil, fmt.Errorf("token expired")
	}

	// Получаем auth по user_id — там хранится секретный ключ (Token)
	auth, err := s.RepoAuth.GetByInt(ctx, "user_id", userID)
	if err != nil {
		return nil, fmt.Errorf("auth record not found for user %d: %w", userID, err)
	}

	// Пересчитываем подпись
	payload := fmt.Sprintf("%d|%d", userID, exp)
	expectedSig := s.Sha256(auth.Secret, payload)
	gotSig, err := base64.RawURLEncoding.DecodeString(string(sigB64))
	if err != nil {
		return nil, fmt.Errorf("invalid signature encoding: %w", err)
	}
	expectedSigBytes, err := base64.RawURLEncoding.DecodeString(expectedSig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode expected signature: %w", err)
	}

	if !hmac.Equal(gotSig, expectedSigBytes) {
		return nil, fmt.Errorf("invalid token signature")
	}

	// Получаем user
	user, err := s.RepoUser.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &domain.UserFull{
		User: user,
		Auth: auth,
	}, nil
}

func (s *Service) UserAuthByPhone(ctx context.Context, phone string) (*domain.UserFull, error) {
	user, err := s.RepoUser.GetBy(ctx, "phone", phone)
	if err != nil {
		return nil, err
	}

	auth, err := s.RepoAuth.GetByInt(ctx, "user_id", user.ID)
	if err != nil {
		return nil, err
	}

	return &domain.UserFull{
		User: user,
		Auth: auth,
	}, nil
}
