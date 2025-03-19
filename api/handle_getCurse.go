package api

import (
	"context"
	"log"
	"net/http"
)

func HandleGetCurse(ctx context.Context, r *http.Request) (*Curse, error) {

	curseId := r.PathValue("curseId")
	p := GetCursePersistence(ctx)

	curse, err := p.Get(ctx, curseId)
	if err != nil {
		log.Println("p.Get:", err)
		return nil, ErrorPersistenceRead
	}
	if curse == nil {
		return nil, ErrorCurseNotFound
	}

	return &curse.Item, nil
}
