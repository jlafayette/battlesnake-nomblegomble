# nomblegomble

Source code for nomblegomble, a [Battlesnake](http://play.battlesnake.com?utm_source=github&utm_medium=readme&utm_campaign=go_starter&utm_content=homepage) written in Go

## Technologies Used

* [Go 1.16](https://golang.org/)

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

Create a python virtual environment and activate it:

```shell
python3 -m venv venv
source venv/bin/activate
```

(On Windows activate the virtual environment using powerscript)

```ps1
.\venv\Script\Activate.ps1
```

Install Python requirements to the virtual environment:

```shell
pip install -r requirements.txt
```

When deploying to a fresh server, run this command to setup the linux service:

```shell
invoke setup-service
```

Run this command to build and deploy to the remote server:

```shell
invoke deploy
```
