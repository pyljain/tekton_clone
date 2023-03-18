package db

import "context"

type Db interface {
	InsertRunner(ctx context.Context, name string) (int, error)
	CreateLink(ctx context.Context, runnerId int, repoName string) error
	FindRunnerByRepository(ctx context.Context, repo string) (int, error)
}
