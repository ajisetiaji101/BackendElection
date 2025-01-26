package handler

import (
	"backend-election/internal/dto"
	"backend-election/internal/pkg/hmac"
	"backend-election/internal/pkg/httpresponse"
	"backend-election/internal/pkg/logger"
	"backend-election/internal/pkg/redis"
	"backend-election/internal/repository"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/julienschmidt/httprouter"
)

// Votes handler
type Votes struct {
	Log   *logger.Logger
	DB    *sql.DB
	Cache *redis.Cache
}

// @Security Bearer
// @Summary Create Votes
// @Description Create Votes
// @Tags Votes
// @Accept  json
// @Produce  json
// @Param Authorization header string true "
// @Success 201 {object} dto.UserResponse
// @Router /votes [post]
func (h *Votes) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = r.Context()

	// Periksa jika context telah dibatalkan atau timeout
	switch ctx.Err() {
	case context.Canceled:
		h.Log.Error(context.Canceled)
		http.Error(w, "Request is canceled", http.StatusRequestTimeout)
		return
	case context.DeadlineExceeded:
		h.Log.Error(context.DeadlineExceeded)
		http.Error(w, "Deadline is exceeded", http.StatusRequestTimeout)
		return
	default:
	}

	// Inisialisasi response struct
	var httpres = httpresponse.Response{Cache: h.Cache}
	// Dekode body permintaan JSON ke dalam VoteRequest
	var voteReq dto.VoteRequest
	defer r.Body.Close()
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&voteReq)
	if err != nil {
		h.Log.Error(err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validasi VoteRequest
	if err := voteReq.Validate(); err != nil {
		h.Log.Error(err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	//validasi kandidat
	var kandidatRepo = repository.KandidatRepository{Log: h.Log, Db: h.DB}

	kandidatRepo.KandidatEntity, err = kandidatRepo.FindByID(ctx, voteReq.CalonID)
	if err != nil {
		http.Error(w, "Kandidat not found", http.StatusNotFound)
		return
	}

	dataKandidatPilih := kandidatRepo.KandidatEntity.GubernurName + " & " + kandidatRepo.KandidatEntity.WakilGubernurName

	//validasi pemilih
	var pemilihRepo = repository.PemilihRepository{Log: h.Log, Db: h.DB}
	pemilih, err := pemilihRepo.FindByID(ctx, voteReq.PemilihID)
	if err != nil {
		http.Error(w, "Pemilih not found", http.StatusNotFound)
		return
	}

	voteReq.CalonID = dataKandidatPilih

	fmt.Printf("Pemilih: %v\n", pemilih.HasVote)

	if pemilih.HasVote {

		fmt.Printf("Pemilih %s sudah melakukan vote", pemilih.PemilihID)

		http.Error(w, "Pemilih sudah melakukan vote", http.StatusConflict)
		return
	}

	// Buat VoteRequest
	requestBody, err := json.Marshal(voteReq)
	if err != nil {
		fmt.Println("Failed to marshal request:", err)
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	// Generate HMAC
	generatedHmac := hmac.GenerateHMAC(os.Getenv("HMAC_KEY_BLOCKCHAIN_ELECTION"), time.Now().Unix())

	fmt.Printf("Generated HMAC: %s\n", generatedHmac)

	// Buat request ke server lain
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8082/vote", bytes.NewBuffer(requestBody))

	fmt.Println("Request to blockchain server : ", req)

	// Set header untuk HMAC
	req.Header.Set("X-HMAC", generatedHmac)
	// Set header untuk timestamp
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))

	if err != nil {
		fmt.Println("Failed to create request:", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Kirim request menggunakan http.Client
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failed to send request to external server:", err)
		http.Error(w, "Failed to send request to external server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Baca respons dari server lain
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Failed to decode response:", err)
		http.Error(w, "Failed to decode response from external server", http.StatusInternalServerError)
		return
	}

	//save pemilih has vote
	pemilih.HasVote = true
	err = pemilihRepo.UpdateVoteUser(ctx, pemilih)
	if err != nil {
		http.Error(w, "Failed to update pemilih", http.StatusInternalServerError)
		return
	}

	// Set response menggunakan httpres.SetMarshal
	httpres.SetMarshal(ctx, w, http.StatusCreated, response, "")
}

// @Summary Show Result
// @Description Show Result
// @Tags Votes
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.UserResponse
// @Router /showresult [get]
func (h *Votes) ShowResult(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = r.Context()

	// Periksa jika context telah dibatalkan atau timeout
	switch ctx.Err() {
	case context.Canceled:
		h.Log.Error(context.Canceled)
		http.Error(w, "Request is canceled", http.StatusRequestTimeout)
		return
	case context.DeadlineExceeded:
		h.Log.Error(context.DeadlineExceeded)
		http.Error(w, "Deadline is exceeded", http.StatusRequestTimeout)
		return
	default:
	}

	// Inisialisasi response struct
	var httpres = httpresponse.Response{Cache: h.Cache}

	//Generate HMAC
	generatedHmac := hmac.GenerateHMAC(os.Getenv("HMAC_KEY_BLOCKCHAIN_ELECTION"), time.Now().Unix())

	fmt.Printf("Generated HMAC: %s\n", generatedHmac)

	// Buat request ke server lain
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8082/showresult", nil)

	fmt.Println("Request to blockchain server : ", req)

	// Set header untuk HMAC
	req.Header.Set("X-HMAC", generatedHmac)
	// Set header untuk timestamp
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))

	if err != nil {
		fmt.Println("Failed to create request:", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Kirim request menggunakan http.Client
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	fmt.Println("Response from blockchain server : ", resp)

	if err != nil {
		fmt.Println("Failed to send request to external server:", err)
		http.Error(w, "Failed to send request to external server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Baca respons dari server lain
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Failed to decode response:", err)
		http.Error(w, "Failed to decode response from external server", http.StatusInternalServerError)
		return
	}

	// Set response menggunakan httpres.SetMarshal
	httpres.SetMarshal(ctx, w, http.StatusOK, response, "")
}
