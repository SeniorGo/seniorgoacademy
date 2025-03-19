package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/SeniorGo/seniorgoacademy/auth"
	"github.com/SeniorGo/seniorgoacademy/persistence"
)

type CreateCurseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func HandleCreateCurse(input *CreateCurseRequest, w http.ResponseWriter, ctx context.Context) (*Curse, error) {

	curse := Curse{
		Id:           uuid.NewString(),
		Author:       auth.GetAuth(ctx).User,
		Title:        input.Title,
		Description:  input.Description,
		CreationTime: time.Now(),
	}
	curse.ModificationTime = curse.CreationTime

	err := curse.Validate()
	if err != nil {
		return nil, err
	}

	p := GetCursePersistence(ctx)
	err = p.Put(context.Background(), &persistence.ItemWithId[Curse]{
		Id:   curse.Id,
		Item: curse,
	})
	if err != nil {
		log.Println("p.Put:", err)
		return nil, ErrorPersistenceWrite
	}

	w.WriteHeader(http.StatusCreated)

	return &curse, nil
}
