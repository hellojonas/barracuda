DEV_BIN := ./bin/barracuda

dev: $(DEV_BIN)
	docker run --rm --env DATABASE_URL=$(DATABASE_URL) barracuda:dev

$(DEV_BIN):
	docker build -t barracuda:dev .
