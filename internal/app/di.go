// Dependency injection / Внедрение зависимостей
package app

type diContainer struct {
	/*	knowApi
		knowService
		knowRepository
		issClient*/
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

/*func (d *diContainer) KnowAPI() api.KnowAPI {
	if d.knowApi == nill {
		d.knowApi = KnowAPI.NewAPI(d.knowService())
	}
}*/
