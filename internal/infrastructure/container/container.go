package container

import (
	"github.com/go-playground/validator/v10"
	"github.com/sarulabs/di"

	"github.com/wojciechpawlinow/usermanagement/internal/application/service"
	"github.com/wojciechpawlinow/usermanagement/internal/config"
	"github.com/wojciechpawlinow/usermanagement/internal/domain/user"
	"github.com/wojciechpawlinow/usermanagement/internal/infrastructure/database/mysql"
	"github.com/wojciechpawlinow/usermanagement/internal/infrastructure/httpserver/handlers"
	"github.com/wojciechpawlinow/usermanagement/pkg/logger"
	"github.com/wojciechpawlinow/usermanagement/pkg/time"
)

func New() di.Container {
	builder, _ := di.NewBuilder()

	if err := builder.Add(di.Def{
		Name: "mysql-conns",
		Build: func(ctn di.Container) (interface{}, error) {
			return mysql.GetConnections(config.Load())
		},
	}); err != nil {
		logger.Fatal(err) // crucial functionality, we can't continue without it
	}

	if err := builder.Add(di.Def{
		Name: "repo-user",
		Build: func(ctn di.Container) (interface{}, error) {
			conns := ctn.Get("mysql-conns").(*mysql.Connections)
			return mysql.NewUserRepository(conns.Read, conns.Write), nil
		},
	}); err != nil {
		logger.Error(err)
	}

	if err := builder.Add(di.Def{
		Name: "service-user",
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewUserService(
				ctn.Get("repo-user").(user.Repository),
				time.NewTimeService(),
			), nil
		},
	}); err != nil {
		logger.Error(err)
	}

	if err := builder.Add(di.Def{
		Name: "http-user",
		Build: func(ctn di.Container) (interface{}, error) {
			return handlers.NewUserHTTPHandler(
				validator.New(),
				ctn.Get("service-user").(service.UserPort),
			), nil
		},
	}); err != nil {
		logger.Error(err)
	}

	return builder.Build()
}
