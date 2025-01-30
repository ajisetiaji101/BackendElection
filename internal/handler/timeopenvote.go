package handler

import (
	"backend-election/internal/model"
	"backend-election/internal/pkg/httpresponse"
	"backend-election/internal/pkg/logger"
	"backend-election/internal/pkg/redis"
	"backend-election/internal/repository"
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type TimeOpenVote struct {
	Log   *logger.Logger
	DB    *sql.DB
	Cache *redis.Cache
}

// @Security Bearer
// @Summary Get time open vote
// @Description Get time open vote
// @Tags TimeOpenVote
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer"
// @Success 200 {object} dto.TimeOpenVoteResponse
// @Router /timeopenvote [get]
func (h *TimeOpenVote) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	var timeOpenVoteRepo = repository.TimeOpenVoteRepository{Log: h.Log, Db: h.DB}
	timeOpenVoteRepo.TimeOpenVoteEntity = model.TimeOpenVote{ID: 1}

	tanggal, err := timeOpenVoteRepo.FindByID(ctx)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var httpres = httpresponse.Response{Cache: h.Cache}

	// Parsing tanggal dari string ke time.Time
	parsedTanggal, err := time.Parse("2006-01-02T15:04:05Z", tanggal)
	if err != nil {
		httpres.SetMarshal(ctx, w, http.StatusInternalServerError, "Internal server error", "")
		return
	}

	// Bandingkan dengan waktu sekarang (UTC)
	if time.Now().UTC().Before(parsedTanggal) {
		httpres.SetMarshal(ctx, w, http.StatusForbidden, "Time open vote failed", "")
		return
	}

	// Jika valid, kirim response sukses
	httpres.SetMarshal(ctx, w, http.StatusOK, "Time open vote success", "")

}
