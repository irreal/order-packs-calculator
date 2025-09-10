# Order Package Calculator Task

This project is an implementation of the take home task when applying to re-partners for a senior go lang position.

I have implemented the solution in steps, from solving the core algorithm, to introducing models, architecture, then implementing the HTTP API, improving DX (hot reload, docker compose), adding persistance using sqlite and finally I added a nice looking web UI. the styling for the UI and some of the content was llm generated.
**AI was used only on some of the presentational stuff to pick fun colors, add style and write copy text. No ai was used for the Go code or the setup of the project at all**

You can check the commit history to trace step by step how the solution was developed.

## Setup Instructions

First step is to clone this repo.

After cloning, copy .env.example into .env and adjust settings as needed

### Running with docker

docker-compose.yml is setup for local development
nothing extra is needed, you can be as simple as `$ docker compose up`

however, for optimal dev experience, consider running:

`$ docker compose up -d && docker compose logs app --follow --tail 20`

this will run the development app in a background container and load up logs so you can see build errors, log output etc.

when you quit logs using `q`, the app continues to run in the background.
Stop the app by running `$ docker compose stop` or `$ docker compose down` to also destroy the containers

### Running locally

You need Go 1.25.1 or greater installed locally.

start the app with:

`$ go run main.go` to try it out.

You can also use air for reloading during development.
Install air if not allready installed:

`$ go install github.com/air-verse/air@latest`

latest should work fine, though if you want to be safe, check the version used in the development docker file at /Dockerfile (1.62.0 as of this writing)

now run:

`$ air`

and off you go, reloading will happen upon any changes in .env or the .go files


## Usage

You can use the app as an API or through the Web UI.

### API

API routes are:
* `GET /api/orders` to get the last 10 orders
* `GET /api/packs` to get the currently used packs
* `POST /api/orders` to create a new order, sample payload: 

```json
{
  "itemCount": 3
}
```

* `POST /api/packs` to change packs in use, sampel payload:

```json
{
  "packs":[1,2,3]
}
```

### Web

On the web, simply navigate to the page and click around.
Start by visting `http://localhost:13131/`
