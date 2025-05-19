PHONY:

compose:
	docker-compose up --build -d
compose-down:
	docker-compose down -v

tests:


migs-up:
	goose up

migs-down:
	goose down

run:
	go build -o ./app cmd/main.go
	./app

bench-post:
	wrk -t8 -c10 -d10s -s post_bench.lua http://0.0.0.0:8080/api/v1/wallet

bench-get:
	wrk -t8 -c10 -d10s -s get_bench.lua http://0.0.0.0:8080/api/v1/wallet