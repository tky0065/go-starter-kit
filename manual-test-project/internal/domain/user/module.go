package user

import "go.uber.org/fx"

var Module = fx.Module("user_service",
	fx.Provide(NewServiceWithJWT),
)
