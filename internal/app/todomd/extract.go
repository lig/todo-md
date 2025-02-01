package todomd

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"runtime"
	"sync"
)

func extractActiveTodos(filenames []string) (todos []string, err error) {
	var wg sync.WaitGroup
	maxGoroutines := runtime.NumCPU()
	semaphore := make(chan struct{}, maxGoroutines)

	todoCh := make(chan []string, len(filenames))
	errCh := make(chan error, len(filenames))

	for _, filename := range filenames {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(filename string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			fileTodos, err := extractTodosFromFile(filename)
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

func extractTodosFromFile(filename string) (entries []string, err error) {
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

		if todo := extractTodoFromLine(line); todo != "" {
			entries = append(entries, fmt.Sprintf("* [%[1]s:%[2]d](%[1]s#L%[2]d): %[3]s", filename, lineNumber, todo))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func extractTodoFromLine(line string) string {
	re := regexp.MustCompile(`(^|\s)TODO:\s*(?P<todo>.*$)`)
	matches := re.FindStringSubmatch(line)

	if len(matches) > 0 {
		return matches[2]
	} else {
		return ""
	}
}
