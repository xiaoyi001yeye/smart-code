package services

import (
    "github.com/kataras/iris/v12"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" 
)


type Repository struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	BranchName string    `db:"branch_name"`
	RepoPath  string    `db:"repo_path"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}


func QueryRepositories(db *sql.DB, name string, branchName string) ([]Repository, error) {
	var repos []Repository
	query := `SELECT id, name, branch_name, repo_path, username, password, created_at, updated_at FROM repositories`

	// 构建 WHERE 子句
	whereClause := ""
	if name != "" {
		whereClause += " AND name = $1"	}	if branchName != "" {		whereClause += " AND branch_name = $2"
	}
	if whereClause != "" {
		query += " WHERE " + whereClause[4:] // 去掉多余的 "AND "
	}

	rows, err := db.Query(query, name, branchName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var repo Repository
		err := rows.StructScan(&repo)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}