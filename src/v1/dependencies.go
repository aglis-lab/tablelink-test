package v1

import (
	"context"
)

type repositories struct {
}

type services struct {
}

type Dependency struct {
	Repositories *repositories
	Services     *services
}

func Dependencies(ctx context.Context) *Dependency {
	repositories := repositories{}

	services := services{}

	return &Dependency{
		Repositories: &repositories,
		Services:     &services,
	}
}
