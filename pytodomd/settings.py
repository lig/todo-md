import pathlib
import typing


TODOMD_DIR: typing.Final = pathlib.Path(__file__).parents[1]

TODOMD_BIN_NAME: typing.Final = 'todo-md'
TODOMD_BIN: typing.Final = TODOMD_DIR / TODOMD_BIN_NAME

TODOMD_UPDATER_NAME: typing.Final = 'todo-md-updater'
TODOMD_UPDATER: typing.Final = TODOMD_DIR / TODOMD_UPDATER_NAME
