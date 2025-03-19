package api

import (
	"log"
	"net/http"
	"sort"
)

func HandleListCurses(w http.ResponseWriter, r *http.Request) ([]Curse, error) {

	ctx := r.Context()

	p := GetCursePersistence(ctx)

	curses, err := p.List(ctx)
	if err != nil {
		log.Println("p.List:", err)
		return nil, ErrorPersistenceRead
	}

	result := make([]Curse, len(curses))
	for i, curse := range curses {
		result[i] = curse.Item
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreationTime.After(result[j].CreationTime)
	})

	return result, nil
}
