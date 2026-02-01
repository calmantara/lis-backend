compose-build:
	docker compose build --no-cache --progress plain

compose-up:
	docker compose up -d

compose-down:
	docker compose down

clean: 
	rm -rf cmd internal tmp go.*