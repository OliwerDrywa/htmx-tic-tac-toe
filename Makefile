.PHONY: tw go 

tw:
	./tailwindcss -i ./web/tailwind.css -o ./web/static/main.css --watch
 

go:
	nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run main.go

	