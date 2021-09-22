import os
import sys
import asyncio
from pathlib import Path

import asyncssh
from dotenv import load_dotenv
from invoke import task


root = Path(__file__).parent.absolute()


load_dotenv()
host = os.environ["REMOTE_HOST"]
port = int(os.environ.get("REMOTE_PORT", 22))
username = os.environ["REMOTE_USER"]
password = os.environ["REMOTE_PASSWORD"]


async def _setup_service():
    async with asyncssh.connect(
        host, port=port, username=username, password=password, known_hosts=None
    ) as conn:
        # setup service
        # copy over .service file and shell script
        service_file = root / "service/battlesnake.service"
        setup_script = root / "service/setup_service.sh"
        async with conn.start_sftp_client() as sftp:
            print("Copying service file to server...", end="")
            await sftp.put(service_file)
            print("done")
            print("Copying setup script to server...", end="")
            await sftp.put(setup_script)
            print("done")

        # run shell script
        print(f"Setting permissions for {setup_script.name}")
        await conn.run(f"chmod +x {setup_script.name}", check=True)
        print(f"Running {setup_script.name}")
        await conn.run(f"./{setup_script.name}")
        # delete shell script
        print(f"Deleting {setup_script.name}")
        await conn.run(f"rm ./{setup_script.name}")


async def _update_battlesnake_server(f):
    async with asyncssh.connect(
        host, port=port, username=username, password=password, known_hosts=None
    ) as conn:
        # copy over shell script and exe
        update_script = root / "service/update.sh"
        async with conn.start_sftp_client() as sftp:
            print(f"Copying {f.name} to server...", end="")
            await sftp.put(f)
            print("done")
            print(f"Copying {update_script.name} to server...", end="")
            await sftp.put(update_script)
            print("done")

        # run shell script
        print(f"Setting permissions for {update_script.name}")
        await conn.run(f"chmod +x {update_script.name}", check=True)
        print(f"Running {update_script.name}")
        await conn.run(f"./{update_script.name}")
        # delete the shell script
        print(f"Deleting {update_script.name}")
        await conn.run(f"rm ./{update_script.name}")


def _build(ctx) -> Path:
    f = "battlesnake-go"
    print("Building")
    go_files = root.glob("*.go")
    go_files_str = " ".join([x.name for x in go_files])
    ctx.run(f"go build -o {f} {go_files_str}", env={"GOOS": "linux", "GOARCH": "amd64"})
    return root / f


@task
def build(ctx):
    _build(ctx)


@task
def setup_service(ctx):
    try:
        asyncio.get_event_loop().run_until_complete(_setup_service())
    except (OSError, asyncssh.Error) as exc:
        sys.exit("SSH connection failed: " + str(exc))


@task
def deploy(ctx):
    f = _build(ctx)
    print("Deploying")
    try:
        asyncio.get_event_loop().run_until_complete(_update_battlesnake_server(f))
    except (OSError, asyncssh.Error) as exc:
        sys.exit("SSH connection failed: " + str(exc))


@task
def format(ctx):
    ctx.run("black -l 100 tasks.py")


# Test Benchmark
# cd tree
# go test -bench Benchmark01 -run xxx > bench01-2.txt