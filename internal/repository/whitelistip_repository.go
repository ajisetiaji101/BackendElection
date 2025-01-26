package repository

import (
	"context"
	"database/sql"

	"backend-election/internal/model"
	"backend-election/internal/pkg/logger"
)

type WhitelistIpRepository struct {
	Db                *sql.DB
	Log               *logger.Logger
	WhitelistIpEntity model.IPWhitelist
}

func (u *WhitelistIpRepository) FindByID(ctx context.Context) error {

	switch ctx.Err() {
	case context.Canceled:
		return u.Log.Error(context.Canceled)
	case context.DeadlineExceeded:
		return u.Log.Error(context.DeadlineExceeded)
	default:
	}

	const q = `SELECT ip_address, created_at FROM ip_whitelist WHERE ip_address = $1`
	// Siapkan prepared statement
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		u.Log.Error(err)
	}
	defer stmt.Close()

	// Eksekusi kueri menggunakan prepared statement
	err = stmt.QueryRowContext(ctx, u.WhitelistIpEntity.IPAddress).Scan(&u.WhitelistIpEntity.IPAddress, &u.WhitelistIpEntity.CreatedAt)

	if err != nil {
		return u.Log.Error(err)
	}

	return nil
}
