# Lightweight Netflix task by Averroes
*Note* there is an .env file with API tokens. Get that from me or things won't work.

Start the app like so:

`make run`

There is a .env.example to look at the secrets that you need

you will need to mount a volume using docker so the data don't go by using the following commands

`mkdir -p $HOME/docker/volumes/postgres`
`docker run --rm   --name pg-docker -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data  postgres`