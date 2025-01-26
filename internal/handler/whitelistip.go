package handler

import (
	"backend-election/internal/dto"
	"backend-election/internal/model"
	"backend-election/internal/pkg/httpresponse"
	"backend-election/internal/pkg/logger"
	"backend-election/internal/pkg/redis"
	"backend-election/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/julienschmidt/httprouter"
)

type WhitelistIp struct {
	Log   *logger.Logger
	DB    *sql.DB
	Cache *redis.Cache
}

// @Security Bearer
// @Summary Get WhitelistIp
// @Description Get WhitelistIp
// @Tags WhitelistIp
// @Accept  json
// @Produce  json
// @Param Authorization header string true "
// @Success 200 {object} dto.UserResponse
// @Router /whitelistip [get]
func (h *WhitelistIp) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

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

	var whitelistipRequest dto.IPWhitelistCreateRequest

	defer r.Body.Close()
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&whitelistipRequest)
	if err != nil {
		h.Log.Error(err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	fmt.Println(whitelistipRequest)

	var whitelistIpRepo = repository.WhitelistIpRepository{Log: h.Log, Db: h.DB}
	whitelistIpRepo.WhitelistIpEntity = model.IPWhitelist{IPAddress: whitelistipRequest.IPAddress}
	err = whitelistIpRepo.FindByID(ctx)
	if err != nil {
		http.Error(w, "Tidak terdaftar", http.StatusInternalServerError)
		return
	}

	var httpres = httpresponse.Response{Cache: h.Cache}
	httpres.SetMarshal(ctx, w, http.StatusOK, whitelistIpRepo.WhitelistIpEntity, "")
}
