package modules

import (
	"proxy/internal/modules/auth/entities"
	"proxy/internal/modules/auth/repository/dbrepo"
	aservice "proxy/internal/modules/auth/service"
	gservice "proxy/internal/modules/geo/service"
	pservice "proxy/internal/modules/proxy/service"
)

type ServicesConfig struct {
	Port      string
	JwtAlg    string
	JwtSecret string
	ApiKey    string
	SecretKey string
	ProxyHost string
	ProxyPort string
}

type Services struct {
	Geo   gservice.GeoServicer
	Auth  aservice.Authenticator
	Proxy pservice.ProxyReverser
}

func NewServices(config ServicesConfig) *Services {
	db := dbrepo.NewMapDBRepo(entities.User{Email: "admin@example.com", Password: "password"})
	return &Services{
		Proxy: pservice.NewProxyReverse(config.ProxyHost, config.ProxyPort),
		Geo:   gservice.NewGeoService(config.ApiKey, config.SecretKey),
		Auth:  aservice.NewUserAuth(config.JwtAlg, config.JwtSecret, db),
	}
}
