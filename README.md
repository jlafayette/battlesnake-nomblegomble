# A Simple [Battlesnake](http://play.battlesnake.com?utm_source=github&utm_medium=readme&utm_campaign=go_starter&utm_content=homepage) Written in Go

## Technologies Used

* [Go 1.16](https://golang.org/)

## TODO List

* Pass value reference if not changing the struct (like state)
* Pass pointer reference if modifying

## Running Tests

```shell
go test
```

## Deploying

Create a server to deploy to with ssh with password enabled (currently using linode).

Create a `.env` file with these variables filled in for your server:

```.env
REMOTE_HOST="xx.xx.xx.xxx"
REMOTE_PORT=xx
REMOTE_USER="username"
REMOTE_PASSWORD="password123"
```

Create a python virtual environment:

```shell
python3 -m venv venv
source venv/bin/activate
```

(On Windows activate using powerscript)

```ps1
.\venv\Script\Activate.ps1
```

Install Python requirements:

```shell
pip install -r requirements.txt
```

Run this command to build and deploy to the remote server:

```shell
invoke deploy
```
