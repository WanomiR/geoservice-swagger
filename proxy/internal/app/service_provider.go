package app

import (
	"proxy/internal/controller"
	"proxy/internal/entities"
	"proxy/internal/repository"
	"proxy/internal/repository/dbrepo"
	"proxy/internal/service"
	"proxy/internal/service/auth"
	"proxy/internal/service/reverse"
	"proxy/utils/readresponder"
)

// TODO: layers initialization

type serviceProvider struct {
	proxyService   service.ProxyReverser
	authService    service.Authenticator
	userController controller.UserController
	userRepository repository.DatabaseRepo
}

func (s *serviceProvider) UserRepository() repository.DatabaseRepo {
	if s.userRepository == nil {
		s.userRepository = dbrepo.NewMapDBRepo(entities.User{Email: "admin@example.com", Password: "password"})
	}
	return s.userRepository
}

func (s *serviceProvider) ProxyService(host, port string) service.ProxyReverser {
	if s.proxyService == nil {
		s.proxyService = reverse.NewProxyReverse(host, port)
	}
	return s.proxyService
}

func (s *serviceProvider) UserController(jwtSecret string) controller.UserController {
	if s.userController == nil {
		readResponder := readresponder.NewReadRespond(readresponder.WithMaxBytes(1 << 20))
		authenticator := auth.NewUserAuth("HS256", jwtSecret)

		s.userController = controller.NewUserControl(
			controller.WithResponder(readResponder),
			controller.WithAuthenticator(authenticator),
		)
	}
	return s.userController
}
