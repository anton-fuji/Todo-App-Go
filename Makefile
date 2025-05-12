DC := docker compose

.PHONY: up upd down build logs ps exec

up:
	$(DC) up --build

upd:
	$(DC) up -d --build

down:
	$(DC) down

build:
	$(DC) build

logs:
	$(DC) logs -f

ps:
	$(DC) ps -a

exec:
	docker exec -it mysql bash