.PHONY:

stock-trader-up:
	docker compose -p stock-trader -f portfolio-service/docker-compose.yaml \
	-f broker-service/docker-compose.yaml \
	up --build

stock-trader-down:
	docker compose -p stock-trader -f portfolio-service/docker-compose.yaml \
	-f broker-service/docker-compose.yaml \
	down

generate-portfolio-migration:
	atlas migrate diff $(migration_name) --dir file://portfolio-service/migrations \
	--to file://portfolio-service/schema.hcl \
	--dev-url docker://mysql/8

generate-broker-migration:
	atlas migrate diff $(migration_name) --dir file://broker-service/migrations \
	--to file://broker-service/schema.hcl \
	--dev-url docker://mysql/8