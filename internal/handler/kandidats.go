package handler

import (
	"backend-election/internal/dto"
	"backend-election/internal/pkg/httpresponse"
	"backend-election/internal/pkg/logger"
	"backend-election/internal/pkg/redis"
	"backend-election/internal/repository"
	"context"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Users handler
type Kandidats struct {
	Log   *logger.Logger
	DB    *sql.DB
	Cache *redis.Cache
}

// @Security Bearer
// @Summary List Users
// @Description List Users
// @Tags Users
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.UserResponse
// @Router /users [get]
func (h *Kandidats) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()

	switch ctx.Err() {
	case context.Canceled:
		h.Log.Error(context.Canceled)
		http.Error(w, "Request is canceled", http.StatusExpectationFailed)
		return
	case context.DeadlineExceeded:
		h.Log.Error(context.DeadlineExceeded)
		http.Error(w, "Deadline is exceeded", http.StatusExpectationFailed)
		return
	default:
	}

	var httpres = httpresponse.Response{Cache: h.Cache}
	var kandidatRepo = repository.KandidatRepository{Log: h.Log, Db: h.DB}
	kandidats, err := kandidatRepo.FindAll(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var kandidatResponse dto.KandidatResponse
	response := kandidatResponse.ListFromEntity(kandidats)
	httpres.SetMarshal(ctx, w, http.StatusOK, response, "")
}
