package runner

import (
	"context"
	"tektonclone/pkg/db"
	"tektonclone/pkg/signing"
)

type Runner struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

func (r *Runner) CreateToken(ctx context.Context, conn db.Db) (string, error) {

	// Persist record
	runnerId, err := conn.InsertRunner(ctx, r.Name)
	if err != nil {
		return "", err
	}

	// Generate new  JWT token
	r.Token, err = signing.GenerateToken(runnerId)
	if err != nil {
		return "", err
	}

	// Return
	return r.Token, nil
}

// r := runner.New("test")
// r.Create() -> 3934797
