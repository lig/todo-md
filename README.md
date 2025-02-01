# TODO.md Pre-Commit Hook

`todo-md` is a pre-commit hook written in Bash that automatically maintains a `TODO.md` file in your repository. It collects `TODO:` comments from your code and organizes them into a markdown file, making it easy to track tasks and improvements.

## Features

- Automatically scans your staged files for `TODO:` comments.
- Updates `TODO.md` with references to the files, line numbers, and corresponding comments.
- Removes outdated entries from `TODO.md` when tasks are removed from the code.
- Uses Forge (e.g. Codeberg or GitHub) style `#L<line number>` links that work in a Forge UI.
- Uses IDE style `:<line number>` text for links that are being recognized in terminal output by most of IDEs.

## Usage

### Installation

1. Install [pre-commit](https://pre-commit.com/).

2. Add the following to your `.pre-commit-config.yaml`:

    ```yaml
    repos:
        - repo: https://codeberg.org/lig/todo-md.git
          rev: v2.0
          hooks:
              - id: todo-md
    ```

3. Install the hook by running:

    ```bash
    pre-commit install
    ```

### How It Works

- The hook scans all staged files for lines containing `TODO:` comments.
- For each `TODO:` comment, it extracts:
  - The file name
  - The line number
  - The comment text
- It updates `TODO.md` in the root of the repository with entries in the format:

    ```markdown
    * [path/to/file:<line_number>](path/to/file#L<line_number>): The TODO comment text
    ```

- Outdated entries (corresponding to removed `TODO:` comments) are removed automatically.

### Example

Given the following code in a staged file:

```python
...
# TODO: Refactor this function
def example():
    pass
```

`TODO.md` will be updated as:

```markdown
* [path/to/file:2](path/to/file#L2): Refactor this function
```

See [TODO.md in this repo](TODO.md) for the real one. :)

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests to the project repository: [todo-md on Codeberg](https://codeberg.org/lig/todo-md).

## Support & Community

Feel free to ask a question or request a feature in [Issues](https://codeberg.org/lig/todo-md/issues) or join our [chat on Matrix](https://matrix.to/#/#todomd:dabar.chat).

## License

This project is licensed under the [MIT License](LICENSE).
