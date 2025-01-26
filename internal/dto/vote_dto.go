package dto

import "errors"

type VoteRequest struct {
	PemilihID string `json:"voter_id"`
	CalonID   string `json:"candidate_id"`
}

func (u *VoteRequest) Validate() error {
	if len(u.CalonID) == 0 {
		return errors.New("Calon ID is required")
	}

	if len(u.PemilihID) == 0 {
		return errors.New("Pemilih ID is required")
	}

	return nil
}
