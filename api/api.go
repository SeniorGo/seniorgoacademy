package api

import (
	"context"
	"net/http"

	"github.com/fulldump/box"

	"github.com/SeniorGo/seniorgoacademy/auth"
	"github.com/SeniorGo/seniorgoacademy/persistence"
	"github.com/SeniorGo/seniorgoacademy/statics"
)

func NewApi(
	version, staticsDir string,
	cursePersistence persistence.Persistencer[Curse],
) http.Handler {

	b := box.NewBox()

	b.WithInterceptors(
		InterceptorAccessLog,
		PrettyError,
	)

	b.WithInterceptors(func(next box.H) box.H {
		return func(ctx context.Context) {
			ctx = context.WithValue(ctx, "curse-persistence", cursePersistence)
			next(ctx)
		}
	})

	b.HandleResourceNotFound = HandleNotFound
	b.HandleMethodNotAllowed = HandleMethodNotAllowed

	b.Handle("GET", "/", HandleRenderHome)
	b.Handle("GET", "/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(version))
	}).WithName("version")

	v0 := b.Group("/v0").WithInterceptors(auth.Require)
	v0.Handle("GET", "/curses", HandleListCurses)
	v0.Handle("POST", "/curses", HandleCreateCurse)
	v0.Handle("GET", "/curses/{curseId}", HandleGetCurse)
	v0.Handle("PATCH", "/curses/{curseId}", HandleModifyCurse)
	v0.Handle("DELETE", "/curses/{curseId}", HandleDeleteCurse)

	// openapi
	buildOpenApi(b)

	// Mount statics
	b.Handle("GET", "/*", statics.ServeStatics(staticsDir)).WithName("serveStatics")

	return b
}
