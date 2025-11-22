package postgres

import (
	"app/internal/app/core/domain"
	"app/pkg/database"
)

type RepoUser struct {
	db database.IDB
}

func NewRepoUser(db database.IDB) *RepoUser {
	return &RepoUser{db: db}
}



func (r *RepoUser) Add(user *domain.User) (int64, error) {
return 0, nil
}


func (r *RepoUser) Get(id int64) (*domain.User, error) {
	return nil, nil
}
func (r *RepoUser) GetBy(filterKey, filterValue string) (*domain.User, error) {
	return nil, nil
}

func (r *RepoUser) List() (*[]domain.User, error) {
	return nil, nil
}
func (r *RepoUser) ListBy(filterKey, filterValue string) (*[]domain.User, error) {
	return nil, nil
}

func (r *RepoUser) Update(id int64, entity *domain.User) error {
	return nil
}
func (r *RepoUser) UpdateBy(filterKey, filterValue string, entity *domain.User) error {
	return nil
}

func (r *RepoUser) Delete(id int64) error {
	return nil
}
func (r *RepoUser) DeleteBy(filterKey, filterValue string) error {
	return nil
}


