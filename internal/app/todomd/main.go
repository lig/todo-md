package todomd

import (
	"os"

	"github.com/alexflint/go-arg"
)

func Run() error {
	arg.MustParse(&args)
	filenames := cleanFilenames(args.Filenames)
	deletedFiles, err := getDeletedFiles()
	if err != nil {
		return err
	}
	filenames = append(filenames, deletedFiles...)

	todoFile, err := os.OpenFile(TODO_FILE, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer todoFile.Close()

	todos, err := loadTodosToKeep(todoFile, filenames)
	if err != nil {
		return err
	}

	activeTodos, err := extractActiveTodos(filenames)
	if err != nil {
		return err
	}
	todos = append(todos, activeTodos...)

	err = writeTodosToFile(todoFile, todos)
	if err != nil {
		return err
	}

	return nil
}
