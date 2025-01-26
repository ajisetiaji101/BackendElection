package model

type Kandidat struct {
	KandidatId        string `db:"kandidat_id"`
	GubernurName      string `db:"gubernur_name"`
	WakilGubernurName string `db:"wakil_gubernur_name"`
	Partai            string `db:"partai"`
	Visi              string `db:"visi"`
	Foto              string `db:"foto"`
	FotoPartai        string `db:"foto_partai"`
}
