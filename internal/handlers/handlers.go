package handlers

type Handlers interface {
	BuildGroups()
	BuildAuthMiddleware()
	BuildRoutes()
	Run(runAddress string) error
}
