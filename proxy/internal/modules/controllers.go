package modules

import (
	acontroller "proxy/internal/modules/auth/controller"
	gcontroller "proxy/internal/modules/geo/controller"
	"proxy/internal/utils/readresponder"
)

type Controllers struct {
	Geo  gcontroller.GeoServicer
	Auth acontroller.Authenticator
}

func NewControllers(services *Services, responder readresponder.ReadResponder) *Controllers {
	return &Controllers{
		Auth: acontroller.NewAuth(services.Auth, responder),
		Geo:  gcontroller.NewGeo(services.Geo, responder),
	}
}
