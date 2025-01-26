package repository

import (
	"context"
	"database/sql"

	"backend-election/internal/model"
	"backend-election/internal/pkg/logger"
)

type KandidatRepository struct {
	Db             *sql.DB
	Log            *logger.Logger
	KandidatEntity model.Kandidat
}

func (u *KandidatRepository) FindAll(ctx context.Context) ([]model.Kandidat, error) {
	const q = `SELECT kandidat_id, gubernur_name, wakil_gubernur_name, foto, foto_partai, partai, visi FROM kandidats`

	// Siapkan prepared statement
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}
	defer stmt.Close()

	// Eksekusi kueri menggunakan prepared statement
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}
	defer rows.Close()

	// Siapkan slice untuk hasil
	var KandidatList []model.Kandidat

	// Iterasi melalui hasil kueri
	for rows.Next() {
		var kandidat model.Kandidat
		if err := rows.Scan(&kandidat.KandidatId, &kandidat.GubernurName, &kandidat.WakilGubernurName, &kandidat.Foto, &kandidat.FotoPartai, &kandidat.Partai, &kandidat.Visi); err != nil {
			u.Log.Error(err)
			return nil, err
		}
		KandidatList = append(KandidatList, kandidat)
	}

	// Periksa error dari iterasi
	if err := rows.Err(); err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return KandidatList, nil
}

func (u *KandidatRepository) FindByID(ctx context.Context, id string) (model.Kandidat, error) {
	const q = `SELECT kandidat_id, gubernur_name, wakil_gubernur_name, foto, foto_partai, partai, visi FROM kandidats WHERE kandidat_id = $1`

	// Siapkan prepared statement
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		u.Log.Error(err)
		return model.Kandidat{}, err
	}
	defer stmt.Close()

	// Eksekusi kueri menggunakan prepared statement
	row := stmt.QueryRowContext(ctx, id)

	// Siapkan variabel untuk hasil
	var kandidat model.Kandidat

	// Ambil hasil kueri
	if err := row.Scan(&kandidat.KandidatId, &kandidat.GubernurName, &kandidat.WakilGubernurName, &kandidat.Foto, &kandidat.FotoPartai, &kandidat.Partai, &kandidat.Visi); err != nil {
		u.Log.Error(err)
		return model.Kandidat{}, err
	}

	return kandidat, nil
}
