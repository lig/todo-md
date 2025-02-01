package todomd

import (
	"errors"

	"github.com/go-git/go-git/v5"
)

func getDeletedFiles() (deletedFiles []string, err error) {
	repo, err := git.PlainOpen(".")
	switch {
	case errors.Is(err, git.ErrRepositoryNotExists):
		return deletedFiles, nil
	case err != nil:
		return nil, err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	status, err := wt.Status()
	if err != nil {
		return nil, err
	}

	for file, s := range status {
		if s.Staging == git.Deleted {
			deletedFiles = append(deletedFiles, file)
		}
	}

	return deletedFiles, nil
}
