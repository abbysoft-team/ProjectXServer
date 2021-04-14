# Gardarike Online server
Gardarike Online is a new MMO strategy and rpg game. Currently it is in early stage of development. We have 2 developers working on the project, one at server side and one making frontend. What we want to achieve is to make gameplay unique but consistent between multiple platforms such as mobile devices and PC's. If you play on your mobile device, the game will be more like strategy economic building simulator. But once you switch platform to PC you will be able to enjoy plain old MMORPG genre expirience.

[![CircleCI](https://circleci.com/gh/abbysoft-team/GardarikeOnlineServer.svg?style=svg)](https://app.circleci.com/pipelines/github/abbysoft-team/GardarikeOnlineServer)

You can find client code here: https://github.com/abbysoft-team/GardarikeOnlineClient 

## Install

You can get fresh release from the Circle CI pipeline. Circle CI is building for ubuntu 16.04 right now.

## Building

Go with modules support is required.

## Linux (Elementary OS Hera) 

Install dependencies
```
sudo apt-get install golang-goprotobuf-dev
sudo apt-get install protobuf-compiler
```

Run in project folder
```
make generate
go build .
```

## Setting up development environment

1. Copy `configs/config.example.toml` file to `configs/config.toml`.
2. Fill `configs/config.toml` [db] section with postgres database connection info
3. Set up local postgres instance
4. Apply all migrations from `db/migrations/` using [go-migrate](https://github.com/golang-migrate/migrate)
5. Now you can run the server!

### Setting up postgres instance

Example commands will be shown for Ubuntu 20.04 disto. If you use some other distro look for it's documentation.

First, download postgres package
```
sudo apt install postgresql
```

Then, start postgres:
```
sudo systemctl start postgresql
```

Make sure postgres up and running:
```
sudo systemctl status postgresql
```

#### Starting postgresql on WSL
You can't use systemctl on WSL, so instead you need to execute:
```
sudo /etc/init.d/postgresql start
```

Then create the database that will be used for the game server. We will assume you've created a new empty database named `gardarike`.

### Applying database migrations
We are using `go-migrate` tool for writing db migrations. You should always maintain your local db in the actual state applying all migrations up to the latest.
First of all, get `go-migrate` tool, how you can obtain `go-migrate` is described on this page: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate.

Apply all migrations in `db/migrations`:
```
migrate -path=db/migrations -database=postgres://user:password@localhost/gardarike up
```

Don't forget to change `user` and `password` to your actuall database user and password.

If all is ok you should see something like that:
```
8/u rotations (27.4826ms) 
```

That means that all migrations up to number 8 was applied to the database `gardarike`. You can check this using this command:
```
migrate -path=db/migrations -database=postgres://admin:admin@localhost/gardarike version
```
You should see the current database version. The version should equals the last migration number in `db/migrations`.

That's all. After all these steps you should have the server and database configured properly and can start contribute to GardarikeOnline!

## LICENSE NOTICE
Feel free to use this code for non-profit goals. If you wan't to use it as part of commercial product contact us via contact@abbysoft.org. Usage without our (maintainers of this repo) permission is prohibited.
