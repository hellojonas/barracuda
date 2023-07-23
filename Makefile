DEV_BIN := ./bin/barracuda

dev: $(DEV_BIN)
	docker run --rm barracuda:dev

$(DEV_BIN):
	docker build -t barracuda:dev .
