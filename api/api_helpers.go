package api

import (
	"context"

	"github.com/SeniorGo/seniorgoacademy/persistence"
)

// GetCursePersistence from context (where curses are stored)
func GetCursePersistence(ctx context.Context) persistence.Persistencer[Curse] {
	p, ok := ctx.Value("curse-persistence").(persistence.Persistencer[Curse])
	if !ok {
		panic("Persistence should be in context")
	}

	return p
}
