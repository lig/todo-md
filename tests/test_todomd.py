import os
import pathlib
import shutil
import tempfile
import typing

import pytest
from inline_snapshot import snapshot

from pytodomd import TODOMD_BIN, TODOMD_BIN_NAME, run


class Dirs(typing.NamedTuple):
    a: pathlib.Path
    b: pathlib.Path


@pytest.fixture
def tmp_dirs() -> typing.Iterator[Dirs]:
    tmpdirs = Dirs(pathlib.Path(tempfile.mkdtemp()), pathlib.Path(tempfile.mkdtemp()))
    yield tmpdirs
    shutil.rmtree(tmpdirs.a, ignore_errors=True)
    shutil.rmtree(tmpdirs.b, ignore_errors=True)


@pytest.fixture
def samples_base_dir() -> pathlib.Path:
    return pathlib.Path(__file__).parent.joinpath('samples')


@pytest.mark.parametrize(
    ('case_name', 'expected_a', 'expected_b'),
    [
        pytest.param(
            '001-new-file',
            snapshot([]),
            snapshot(['* [sample.txt:3](sample.txt#L3): todo text 1']),
            id='new file',
        ),
        pytest.param(
            '002-update-file',
            snapshot(['* [sample.txt:3](sample.txt#L3): todo text 1']),
            snapshot(
                ['* [sample.txt:4](sample.txt#L4): updated todo text on another line']
            ),
            id='update file',
        ),
        pytest.param(
            '003-delete-file',
            snapshot(['* [sample.txt:3](sample.txt#L3): todo text 1']),
            snapshot([]),
            id='delete file',
        ),
        pytest.param(
            '004-numeric-ordering',
            snapshot(
                [
                    '* [sample.txt:2](sample.txt#L2): todo text on line 2',
                    '* [sample.txt:10](sample.txt#L10): todo text on line 10',
                ]
            ),
            snapshot([]),
            id='delete file',
        ),
    ],
)
def test_todomd(
    tmp_dirs: Dirs,
    samples_base_dir: pathlib.Path,
    case_name: str,
    expected_a: str,
    expected_b: str,
) -> None:
    samples_dirs = Dirs(
        samples_base_dir / case_name / 'a', samples_base_dir / case_name / 'b'
    )

    # Run A
    files_a = [
        pathlib.Path(file).relative_to(samples_dirs.a)
        for file in samples_dirs.a.glob('**')
        if file.is_file()
    ]

    shutil.copytree(samples_dirs.a, tmp_dirs.a, dirs_exist_ok=True)
    todomd_updater = tmp_dirs.a / TODOMD_BIN_NAME
    shutil.copy(TODOMD_BIN, todomd_updater)

    os.chdir(tmp_dirs.a)
    run(todomd_bin=todomd_updater, files=files_a)

    assert tmp_dirs.a.joinpath('TODO.md').read_text().splitlines() == expected_a

    # Copy `TODO.md` between runs`
    shutil.copy(tmp_dirs.a / 'TODO.md', tmp_dirs.b / 'TODO.md')

    # Run B
    files_b = files_a + [
        pathlib.Path(file).relative_to(samples_dirs.b)
        for file in samples_dirs.b.glob('**')
        if file.is_file()
    ]

    shutil.copytree(samples_dirs.b, tmp_dirs.b, dirs_exist_ok=True)
    todomd_updater = tmp_dirs.b / TODOMD_BIN_NAME
    shutil.copy(TODOMD_BIN, todomd_updater)

    os.chdir(tmp_dirs.b)
    run(todomd_bin=todomd_updater, files=files_b)

    assert tmp_dirs.b.joinpath('TODO.md').read_text().splitlines() == expected_b
