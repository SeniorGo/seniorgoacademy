package api

import (
	"log"
	"net/http"
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
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

	collator := collate.New(language.Spanish, collate.Numeric)

	sort.Slice(result, func(i, j int) bool {
		return -1 == collator.CompareString(result[i].Title, result[j].Title)
	})

	return result, nil
}
