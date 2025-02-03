package todomd

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"regexp"
	"runtime"
	"sync"
)

var sourceRegex = regexp.MustCompile(`(^|\s)TODO:\s*(?P<todo>.*$)`)

func extractActiveTodos(filenames []string) (todos []*Todo, err error) {
	var wg sync.WaitGroup
	maxGoroutines := runtime.NumCPU()
	semaphore := make(chan struct{}, maxGoroutines)

	todoCh := make(chan []*Todo, len(filenames))
	errCh := make(chan error, len(filenames))

	for _, filename := range filenames {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(filename string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			fileTodos, err := extractTodosFromSourceFile(filename)
			if err != nil {
				errCh <- err
				return
			}
			todoCh <- fileTodos
		}(filename)
	}

	go func() {
		wg.Wait()
		close(todoCh)
		close(errCh)
	}()

	for fileTodos := range todoCh {
		todos = append(todos, fileTodos...)
	}

	if len(errCh) > 0 {
		return nil, <-errCh
	}

	return todos, nil
}

func extractTodosFromSourceFile(filename string) (entries []*Todo, err error) {
	file, err := os.Open(filename)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return entries, nil
	case err != nil:
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		if text := extractTextFromSourceLine(line); text != "" {
			entries = append(entries, &Todo{
				Filename:   filename,
				LineNumber: lineNumber,
				Text:       text,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func extractTextFromSourceLine(line string) string {
	matches := sourceRegex.FindStringSubmatch(line)

	if len(matches) > 0 {
		return matches[2]
	} else {
		return ""
	}
}
