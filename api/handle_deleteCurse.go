package api

import (
	"context"
	"log"
	"net/http"

	"github.com/SeniorGo/seniorgoacademy/auth"
)

func HandleDeleteCurse(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	curseId := r.PathValue("curseId")
	p := GetCursePersistence(ctx)

	curse, err := p.Get(ctx, curseId)
	if err != nil {
		log.Println("p.Get", err)
		return ErrorPersistenceRead
	}
	if curse == nil {
		return ErrorCurseNotFound
	}

	// Access control
	if curse.Item.Author.ID != auth.GetAuth(ctx).User.ID {
		return ErrorCurseForbidden
	}

	err = p.Delete(ctx, curseId)
	if err != nil {
		log.Println("p.Delete:", err)
		return ErrorPersistenceWrite
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
