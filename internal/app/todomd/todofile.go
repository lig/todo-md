package todomd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

const TODO_FILE = "./TODO.md"
const TODO_TEMPLATE = "* [{{.Filename}}:{{.LineNumber}}]({{.Filename}}#L{{.LineNumber}}): {{.Text}}\n"

var todoRegex = regexp.MustCompile(`^* \[.*?\]\((?P<filename>.+?)#L(?P<lineNumber>\d+)\): (?P<text>.*)$`)

func loadTodosToKeep(todoFile *os.File, filenames []string) (todos []*Todo, err error) {
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
			todo, err := extractTodoFromTodoLine(line)
			if err != nil {
				return nil, err
			}
			todos = append(todos, todo)
		}
	}

	return todos, scanner.Err()
}

func writeTodosToFile(todoFile *os.File, todos []*Todo) error {
	todoTemplate, err := template.New("todoTemplate").Parse(TODO_TEMPLATE)
	if err != nil {
		return err
	}

	if err := todoFile.Truncate(0); err != nil {
		return err
	}
	if _, err := todoFile.Seek(0, 0); err != nil {
		return err
	}

	sort.Slice(todos, func(i, j int) bool {
		a, b := todos[i], todos[j]
		return a.Filename < b.Filename && a.LineNumber < b.LineNumber
	})

	for _, todo := range todos {
		err = todoTemplate.Execute(todoFile, todo)
		if err != nil {
			return err
		}
	}

	return err
}

func extractTodoFromTodoLine(line string) (*Todo, error) {
	matches := todoRegex.FindStringSubmatch(line)

	filename, lineNumberStr, text := matches[1], matches[2], matches[3]
	lineNumber, err := strconv.Atoi(lineNumberStr)
	if err != nil {
		return nil, err
	}
	return &Todo{Filename: filename, LineNumber: lineNumber, Text: text}, nil
}
