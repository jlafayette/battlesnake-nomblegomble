import os
import sys
import asyncio

import asyncssh
from dotenv import load_dotenv
from invoke import task



async def copy_file(path):
    load_dotenv()
    host = os.environ["HOST"]
    port = int(os.environ["PORT"])
    username = os.environ["REMOTE_USER"]
    password = os.environ["REMOTE_PASSWORD"]
    async with asyncssh.connect(host,
                                port=port,
                                username=username,
                                password=password,
                                known_hosts=None) as conn:
        async with conn.start_sftp_client() as sftp:
            is_file = await sftp.isfile(path)
            if is_file:
                await sftp.remove(path)
            await sftp.put(path)
        await conn.run(f"chmod +x {path}", check=True)


@task
def deploy(ctx):
    f = "main"
    ctx.run(f"go build -o {f} main.go logic.go", env={"GOOS": "linux", "GOARCH": "amd64"})
    print("Deploying...")
    try:
        asyncio.get_event_loop().run_until_complete(copy_file(f))
    except (OSError, asyncssh.Error) as exc:
        sys.exit('SSH connection failed: ' + str(exc))
