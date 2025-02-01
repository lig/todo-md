import pathlib
import typing


TODOMD_DIR: typing.Final = pathlib.Path(__file__).parents[1] / 'build'

TODOMD_BIN_NAME: typing.Final = 'todo-md'
TODOMD_BIN: typing.Final = TODOMD_DIR / TODOMD_BIN_NAME
