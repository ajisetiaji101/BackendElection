package dto

import "backend-election/internal/model"

type KandidatResponse struct {
	KandidatId        string `db:"kandidat_id"`
	GubernurName      string `db:"gubernur_name"`
	WakilGubernurName string `db:"wakil_gubernur_name"`
	Partai            string `db:"partai"`
	Visi              string `db:"visi"`
	Foto              string `db:"foto"`
	FotoPartai        string `db:"foto_partai"`
}

func (u *KandidatResponse) ListFromEntity(kandidats []model.Kandidat) []KandidatResponse {
	var kandidatResponses []KandidatResponse
	for _, kandidat := range kandidats {
		kandidatResponses = append(kandidatResponses, KandidatResponse{
			KandidatId:        kandidat.KandidatId,
			GubernurName:      kandidat.GubernurName,
			WakilGubernurName: kandidat.WakilGubernurName,
			Partai:            kandidat.Partai,
			Visi:              kandidat.Visi,
			Foto:              kandidat.Foto,
			FotoPartai:        kandidat.FotoPartai,
		})
	}
	return kandidatResponses
}
