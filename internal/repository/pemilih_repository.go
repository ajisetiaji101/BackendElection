package repository

import (
	"backend-election/internal/model"
	"backend-election/internal/pkg/logger"
	"context"
	"database/sql"
)

type PemilihRepository struct {
	Db            *sql.DB
	Log           *logger.Logger
	PemilihEntity model.Pemilih
}

func (u *PemilihRepository) FindAll(ctx context.Context) ([]model.Pemilih, error) {
	const q = `SELECT pemilih_id, nik, nama, tanggal_lahir, alamat, sidik_jari FROM pemilih`

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
	var pemilihList []model.Pemilih

	// Iterasi melalui hasil kueri
	for rows.Next() {
		var pemilih model.Pemilih
		if err := rows.Scan(&pemilih.PemilihID, &pemilih.NIK, &pemilih.Nama, &pemilih.TanggalLahir, &pemilih.Alamat, &pemilih.SidikJari); err != nil {
			u.Log.Error(err)
			return nil, err
		}
		pemilihList = append(pemilihList, pemilih)
	}

	// Periksa error dari iterasi
	if err := rows.Err(); err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return pemilihList, nil
}

func (u *PemilihRepository) Find(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return u.Log.Error(context.Canceled)
	case context.DeadlineExceeded:
		return u.Log.Error(context.DeadlineExceeded)
	default:
	}

	const q = `SELECT pemilih_id, nik, nama, tanggal_lahir, alamat, sidik_jari FROM pemilih WHERE pemilih_id=$1`
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		return u.Log.Error(err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.PemilihEntity.PemilihID).Scan(&u.PemilihEntity.PemilihID, &u.PemilihEntity.NIK, &u.PemilihEntity.Nama, &u.PemilihEntity.TanggalLahir, &u.PemilihEntity.Alamat, &u.PemilihEntity.SidikJari)
	if err != nil {
		return u.Log.Error(err)
	}
	return nil
}

func (u *PemilihRepository) FindByID(ctx context.Context, id string) (model.Pemilih, error) {
	const q = `SELECT pemilih_id, nik, nama, tanggal_lahir, alamat, sidik_jari, has_vote, tps FROM pemilih WHERE pemilih_id = $1`

	// Siapkan prepared statement
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		u.Log.Error(err)
		return model.Pemilih{}, err
	}
	defer stmt.Close()

	// Eksekusi kueri menggunakan prepared statement
	row := stmt.QueryRowContext(ctx, id)

	// Siapkan variabel untuk hasil
	var pemilih model.Pemilih

	// Ambil hasil kueri
	if err := row.Scan(&pemilih.PemilihID, &pemilih.NIK, &pemilih.Nama, &pemilih.TanggalLahir, &pemilih.Alamat, &pemilih.SidikJari, &pemilih.HasVote, &pemilih.TPS); err != nil {
		u.Log.Error(err)
		return model.Pemilih{}, err
	}

	return pemilih, nil
}

func (u *PemilihRepository) FindByNIK(ctx context.Context, id string) (model.Pemilih, error) {

	const q = `SELECT pemilih_id, nik, nama, tanggal_lahir, alamat, sidik_jari, has_vote, tps FROM pemilih WHERE nik=$1`
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		u.Log.Error(err)
		return model.Pemilih{}, err
	}
	defer stmt.Close()

	// Eksekusi kueri menggunakan prepared statement
	row := stmt.QueryRowContext(ctx, id)

	// Siapkan variabel untuk hasil
	var pemilih model.Pemilih

	// Ambil hasil kueri
	if err := row.Scan(&pemilih.PemilihID, &pemilih.NIK, &pemilih.Nama, &pemilih.TanggalLahir, &pemilih.Alamat, &pemilih.SidikJari, &pemilih.HasVote, &pemilih.TPS); err != nil {
		u.Log.Error(err)
		return model.Pemilih{}, err
	}

	return pemilih, nil
}

func (u *PemilihRepository) UpdateVoteUser(ctx context.Context, pemilih model.Pemilih) error {
	const q = `UPDATE pemilih SET has_vote = $1 WHERE pemilih_id = $2`

	// Siapkan prepared statement
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		u.Log.Error(err)
		return err
	}
	defer stmt.Close()

	// Eksekusi kueri menggunakan prepared statement
	_, err = stmt.ExecContext(ctx, pemilih.HasVote, pemilih.PemilihID)
	if err != nil {
		u.Log.Error(err)
		return err
	}

	return nil
}

func (u *PemilihRepository) GetTotalPemilihPersenByTpsAndVote(ctx context.Context) (int, error) {
	const qTotalHasVote = `SELECT COUNT(pemilih_id) FROM pemilih WHERE has_vote = $1 and tps = $2`
	const qTotalPemilih = `SELECT COUNT(pemilih_id) FROM pemilih WHERE tps = $1`

	// Siapkan prepared statement
	stmtTotalHasVote, err := u.Db.PrepareContext(ctx, qTotalHasVote)
	if err != nil {
		u.Log.Error(err)
		return 0, err
	}
	defer stmtTotalHasVote.Close()

	stmtTotalPemilih, err := u.Db.PrepareContext(ctx, qTotalPemilih)
	if err != nil {
		u.Log.Error(err)
		return 0, err
	}
	defer stmtTotalPemilih.Close()

	// Eksekusi kueri menggunakan prepared statement
	var totalHasVote, totalPemilih int
	err = stmtTotalHasVote.QueryRowContext(ctx, true, u.PemilihEntity.TPS).Scan(&totalHasVote)
	if err != nil {
		u.Log.Error(err)
		return 0, err
	}

	err = stmtTotalPemilih.QueryRowContext(ctx, u.PemilihEntity.TPS).Scan(&totalPemilih)
	if err != nil {
		u.Log.Error(err)
		return 0, err
	}

	if totalPemilih == 0 {
		return 0, nil
	}

	return (totalHasVote * 100) / totalPemilih, nil
}
