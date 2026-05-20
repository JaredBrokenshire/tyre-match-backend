build:
	docker compose up --build -d
test:
	docker compose -f docker-compose-test.yml run --rm tyre_match_test
test-coverage:
	docker compose -f docker-compose-test.yml run --rm tyre_match_test pytest --cov=. --cov-report=term-missing --cov-report=html:/main/test-artifacts/coverage