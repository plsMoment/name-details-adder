run:
	docker-compose up --build name-details-adder

down:
	docker compose down name-details-adder

.PHONY: run down