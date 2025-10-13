package usecase

import (
	"context"
	"fmt"
	"time"
)

var ErrDatabaseTimeout error = fmt.Errorf("время ожидания базы данных истекло")

type UserService struct {
	User       User
	repository UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repository: repo,
	}
}

func (u *UserService) UserByUid(uid string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := u.repository.Read(ctx, uid)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("ctx deadline %w: %v", ErrDatabaseTimeout, err)
		}
		return nil, fmt.Errorf("failed to read user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	u.User = *user
	return user, nil
}

func (u *UserService) CreateUser(user *User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uid, err := u.repository.Create(ctx, user)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("ctx deadline %w: %v", ErrDatabaseTimeout, err)
		}
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	u.User = *user
	u.User.Uid = uid
	return uid, nil
}

func (u *UserService) UpdateUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := u.repository.Update(ctx, user); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("ctx deadline %w: %v", ErrDatabaseTimeout, err)
		}
		return fmt.Errorf("failed to update user: %w", err)
	}
	u.User = *user
	return nil
}

func (u *UserService) DeleteByUid(uid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := u.repository.Delete(ctx, uid); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("ctx deadline %w: %v", ErrDatabaseTimeout, err)
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	u.User = User{}
	return nil
}

func (u *UserService) AllUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := u.repository.All(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("ctx deadline %w: %v", ErrDatabaseTimeout, err)
		}
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}

/*
func (u *UserService) CheckAuth(login, password string) (*User, error) {
    users, err := u.AllUsers()
    if err != nil {
        return nil, err
    }

    for _, user := range users {
        if user.Login == login && user.Password == password {
            return user, nil
        }
    }

    return nil, fmt.Errorf("invalid credentials")
}


func (u *UserService) CheckLogin(login string) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    users, err := u.repository.All(ctx)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return false, fmt.Errorf("ctx deadline %w: %v", ErrDatabaseTimeout, err)
        }
        return false, fmt.Errorf("failed to check login: %w", err)
    }

    for _, user := range users {
        if user.Login == login {
            return true, nil
        }
    }

    return false, nil
}
*/
