build:
	docker-compose build

start: stop build
	docker-compose up -d

stop:
	docker-compose down

clean: stop
	rm -rf quoter-rust/target
	rm -rf fetcher-go/fetcher

reinstall: clean build

restart: stop start
