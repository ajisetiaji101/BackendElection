package dto

import (
	"backend-election/internal/model"
)

type IPWhitelistCreateRequest struct {
	IPAddress string `json:"ip_address"`
}

func (u *IPWhitelistCreateRequest) ToEntity() model.IPWhitelist {
	return model.IPWhitelist{
		IPAddress: u.IPAddress,
	}
}
