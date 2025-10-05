package handlers

type Handlers interface {
	BuildGroups()
	BuildAuthMiddleware(jwtSecret string)
	BuildRoutes()
	Run() error
}
