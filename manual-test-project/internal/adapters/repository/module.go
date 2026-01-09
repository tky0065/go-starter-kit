package repository

import (
	"go.uber.org/fx"
	"manual-test-project/internal/domain/user"
	"manual-test-project/internal/interfaces"
)

var Module = fx.Module("repository",
	fx.Provide(
		fx.Annotate(
			NewUserRepository,
			fx.As(new(interfaces.UserRepository)),
			fx.As(new(user.UserRepository)),
		),
	),
)
