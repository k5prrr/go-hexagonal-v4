package service

import (
	"app/internal/app/core/domain"
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

func (s *Service) hash(secret, data []byte) []byte {
	secretWithSalt := append(secret, s.salt...)

	h := hmac.New(sha256.New, secretWithSalt)
	h.Write(data)

	return h.Sum(nil)
}
func (s *Service) Token(authID int64, secret string) string {
	payload := fmt.Sprintf("%d|%d", authID, time.Now().Add(30*24*time.Hour).Unix())
	payloadByte := []byte(payload)

	tokenRaw := fmt.Sprintf(
		"%s|%s|%s",
		payload,
		s.hash(s.salt, payloadByte),
		s.hash([]byte(secret), payloadByte),
	)

	return base64.RawURLEncoding.EncodeToString([]byte(tokenRaw))
}
func (s *Service) CurrentUser(ctx context.Context, token string) (*domain.UserFull, error) {
	userFull := s.MapAuth.Get(token)
	if userFull != nil {
		return userFull, nil
	}

	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token encoding: %w", err)
	}

	parts := bytes.Split(data, []byte{'|'})
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid token format")
	}

	authIDRaw, expRaw, sigSaltB64, sigSecretB64 := parts[0], parts[1], parts[2], parts[3]

	authID, err := strconv.ParseInt(string(authIDRaw), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid auth id: %w", err)
	}

	exp, err := strconv.ParseInt(string(expRaw), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid expiration time: %w", err)
	}

	if time.Now().Unix() > exp {
		return nil, fmt.Errorf("token expired")
	}

	payload := fmt.Sprintf("%d|%d", authID, exp)
	payloadByte := []byte(payload)

	// Проверка соли
	newSigSalt := s.hash(s.salt, payloadByte)
	oldSigSalt, err := base64.RawURLEncoding.DecodeString(string(sigSaltB64))
	if err != nil {
		return nil, fmt.Errorf("invalid signature encoding: %w", err)
	}
	if !hmac.Equal(oldSigSalt, newSigSalt) {
		return nil, fmt.Errorf("invalid token signature")
	}

	// Получаем auth - там хранится секретный ключ (Token)
	auth, err := s.RepoAuth.Get(ctx, authID)
	if err != nil {
		return nil, fmt.Errorf("auth record not found %d: %w", authID, err)
	}

	// Проверка Секрета
	newSigSecret := s.hash([]byte(auth.Secret), payloadByte)
	oldSigSecret, err := base64.RawURLEncoding.DecodeString(string(sigSecretB64))
	if err != nil {
		return nil, fmt.Errorf("invalid signature encoding: %w", err)
	}
	if !hmac.Equal(oldSigSecret, newSigSecret) {
		return nil, fmt.Errorf("invalid token signature")
	}

	// Получаем user
	user, err := s.RepoUser.Get(ctx, authID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// От его роли получаем остальное

	userFull = &domain.UserFull{
		User: user,
		Auth: auth,
	}
	s.MapAuth.Set(token, userFull)

	return userFull, nil
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
