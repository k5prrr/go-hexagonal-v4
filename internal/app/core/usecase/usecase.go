package usecase

import (
	"app/internal/app/core/domain"
	"app/internal/app/core/port"
	"app/internal/app/core/service"
	"context"
)

// Тут он как основной, тут только создания
type UseCase struct {
	service *service.Service
}

func New(service *service.Service) port.IUseCase {
	return &UseCase{
		service: service,
	}
}

func (u *UseCase) Test1() (*domain.User, error) {
	ctx := context.Background()

	/*user := domain.User{
		Name:       "name4",
		FamilyName: "jo",
		Phone:      fmt.Sprintf("phone%d", rand.Intn(999)),
	}*/

	// Осталось лист листб, дел, делб,
	return u.service.RepoUser.GetBy(ctx, "phone", "phone97")
	//return u.service.RepoUser.Get(ctx, 25)

	//err := u.service.RepoUser.UpdateBy(ctx, "name", "name3", &user)
	//err := u.service.RepoUser.Update(ctx, 29, &user)
	//return 0, err
	//return u.service.RepoUser.Add(ctx, &user)
}
