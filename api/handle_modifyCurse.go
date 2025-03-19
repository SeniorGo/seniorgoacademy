package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/SeniorGo/seniorgoacademy/auth"
)

type ModifyCurseRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func HandleModifyCurse(ctx context.Context, r *http.Request, input *ModifyCurseRequest) (*Curse, error) {

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

	// Access control
	if curse.Item.Author.ID != auth.GetAuth(ctx).User.ID {
		return nil, ErrorCurseForbidden
	}

	curse.Item.Author = auth.GetAuth(ctx).User // Update user data
	curse.Item.ModificationTime = time.Now()

	if input.Title != nil {
		curse.Item.Title = *input.Title
	}

	if input.Description != nil {
		curse.Item.Description = *input.Description
	}

	err = curse.Item.Validate()
	if err != nil {
		return nil, err
	}

	err = p.Put(ctx, curse)
	if err != nil {
		log.Println("p.Put:", err)
		return nil, ErrorPersistenceWrite
	}

	return &curse.Item, nil
}
