package dto

import "backend-election/internal/model"

type OpenVoteTimeResponse struct {
	OpenVoteTime string `json:"date"`
}

type OpenVoteTimeRequest struct {
	OpenVoteTime string `json:"date"`
}

func (u *OpenVoteTimeRequest) ToEntity() model.TimeOpenVote {
	return model.TimeOpenVote{
		OpenVoteTime: u.OpenVoteTime,
	}
}
