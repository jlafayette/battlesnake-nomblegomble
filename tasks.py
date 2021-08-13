import os
import sys
import asyncio

import asyncssh
from dotenv import load_dotenv
from invoke import task


async def copy_file(path):
    load_dotenv()
    host = os.environ["REMOTE_HOST"]
    port = int(os.environ.get("REMOTE_PORT", 22))
    username = os.environ["REMOTE_USER"]
    password = os.environ["REMOTE_PASSWORD"]
    async with asyncssh.connect(
        host, port=port, username=username, password=password, known_hosts=None
    ) as conn:
        kill_result = await conn.run(f"kill -9 `pidof ./{path}`")
        print(f"# kill -9 `pidof ./{path}`\n{kill_result.returncode} {kill_result.stdout}")
        if kill_result.returncode == 0:
            print("Killed currently running server")
        async with conn.start_sftp_client() as sftp:
            is_file = await sftp.isfile(path)
            if is_file:
                print(f"Removing existing file: {path}")
                await sftp.remove(path)
            print("Copying new file to server...", end="")
            await sftp.put(path)
            print("done")
        print("Setting permissions")
        await conn.run(f"chmod +x {path}", check=True)
        print("Starting server...")
        await conn.run(f"setsid ./{path} < /dev/zero 2>&1&> /dev/null &")
        print("Done")


@task
def deploy(ctx):
    f = "main"
    print("Building")
    ctx.run(f"go build -o {f} main.go logic.go grid.go util.go logging.go", env={"GOOS": "linux", "GOARCH": "amd64"})
    print("Deploying")
    try:
        asyncio.get_event_loop().run_until_complete(copy_file(f))
    except (OSError, asyncssh.Error) as exc:
        sys.exit("SSH connection failed: " + str(exc))


@task
def format(ctx):
    ctx.run("black -l 100 tasks.py")
