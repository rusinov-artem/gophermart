.PHONY: bash
bash: 
	docker-compose run --rm gophermart-app bash

.PHONY: test
test: 
	docker-compose run --rm gophermart-app bash test/test.bash

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: down-v
down-v:
	docker-compose down -v

.PHONY: clean
clean:
	docker-compose down -v --remove-orphans --rmi all