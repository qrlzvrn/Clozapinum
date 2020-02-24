conf:
	./initconf.sh
	echo "OK"

conf-clean:
	rm -rf .env/

docker-run:
	docker-compose up --build -d