package handlers

type Handlers interface {
	BuildGroups()
	BuildAuthMiddleware()
	BuildRoutes()
	Run() error
}
