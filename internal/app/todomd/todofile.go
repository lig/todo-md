package todomd

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

const TODO_FILE = "./TODO.md"

func loadTodosToKeep(todoFile *os.File, filenames []string) (todos []string, err error) {
	scanner := bufio.NewScanner(todoFile)
	for scanner.Scan() {
		line := scanner.Text()
		shouldInclude := true

		for _, filename := range filenames {
			if strings.Contains(line, fmt.Sprintf("[%s:", filename)) {
				shouldInclude = false
				break
			}
		}

		if shouldInclude {
			todos = append(todos, line)
		}
	}

	return todos, scanner.Err()
}

func writeTodosToFile(todoFile *os.File, lines []string) error {
	if err := todoFile.Truncate(0); err != nil {
		return err
	}
	if _, err := todoFile.Seek(0, 0); err != nil {
		return err
	}

	slices.Sort(lines)
	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n"
	}
	_, err := todoFile.WriteString(content)
	return err
}
