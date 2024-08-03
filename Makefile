.PHONY: bash
bash: 
	docker-compose run --rm gophermart-app bash

.PHONY: test
test: 
	docker-compose run --rm gophermart-app bash test/test.bash
