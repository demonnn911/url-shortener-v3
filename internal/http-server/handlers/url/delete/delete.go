package delete

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener-3rd-attempt/internal/lib/api/response"
	"url-shortener-3rd-attempt/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("field alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("no url with alias", slog.String("alias", alias))
				render.JSON(w, r, resp.Error("invalid request"))
				return
			} else {
				log.Info("failed to delete url")
				render.JSON(w, r, resp.Error("internal error"))
				return
			}
		}

		log.Info("delete alias", slog.Any("alias", alias))
		render.JSON(w, r, resp.OK())
	}
}
