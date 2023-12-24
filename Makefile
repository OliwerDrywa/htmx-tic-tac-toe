.PHONY: run build 

run:
	./tailwindcss -i ./web/tailwind.css -o ./web/static/main.css --watch &
	go run main.go

build:
	./tailwindcss -i ./web/tailwind.css -o ./web/static/main.css
	go build 
	