package repository

import (
	"context"
	"database/sql"

	"backend-election/internal/model"
	"backend-election/internal/pkg/logger"
)

type TimeOpenVoteRepository struct {
	Db                 *sql.DB
	Log                *logger.Logger
	TimeOpenVoteEntity model.TimeOpenVote
}

func (u *TimeOpenVoteRepository) FindByID(ctx context.Context) (string, error) {

	// Cek jika context mengalami error sebelum eksekusi query
	switch ctx.Err() {
	case context.Canceled:
		u.Log.Error(context.Canceled)
		return "", context.Canceled
	case context.DeadlineExceeded:
		u.Log.Error(context.DeadlineExceeded)
		return "", context.DeadlineExceeded
	}

	const q = `SELECT date FROM time_open_vote WHERE id = $1`

	// Siapkan prepared statement
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		u.Log.Error(err)
		return "", err
	}
	defer stmt.Close()

	// Eksekusi kueri menggunakan prepared statement
	err = stmt.QueryRowContext(ctx, u.TimeOpenVoteEntity.ID).Scan(&u.TimeOpenVoteEntity.OpenVoteTime)
	if err != nil {
		u.Log.Error(err)
		return "", err
	}

	return u.TimeOpenVoteEntity.OpenVoteTime, nil
}
