package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}

	return false
}

type PermissionModel struct {
	DB *sql.DB
}

func (m PermissionModel) ReadAllForUser(userID int64) (Permissions, error) {
	query := `SELECT permission.code
						FROM permission
						INNER JOIN user_permission ON user_permission.permission_id = permission.id
						INNER JOIN app_user ON user_permission.user_id = app_user.id
						WHERE app_user.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (m PermissionModel) AddForUser(userID int64, codes ...string) error {
	query := `INSERT INTO user_permission
						SELECT $1, permission.id FROM permission WHERE permission.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID, pq.Array(codes))
	return err
}
