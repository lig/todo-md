import pathlib
import subprocess

from .settings import TODOMD_BIN


def run(*, todomd_bin: pathlib.Path = TODOMD_BIN, files=list[pathlib.Path]) -> None:
    subprocess.run([todomd_bin, *[str(file) for file in files]])
