package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetPipelineDef(ctx context.Context, repositoryName string, commitRef string) ([]byte, error) {

	loc := fmt.Sprintf("./tmp/%s", commitRef)
	repo, err := git.PlainClone(loc, false, &git.CloneOptions{
		URL:      repositoryName,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commitRef),
	})
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("./tmp/%s/pipeline.yaml", commitRef)
	filebytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return filebytes, nil
}
