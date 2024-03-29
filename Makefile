run: 1 2 3

1:
	@templ generate

2:
	@tailwindcss -i web/main.css -o web/_main.css

3:
	@go run cmd/main.go

watch:
	@make 1w
	@make 2w
	@make 3w

1w:
	@templ generate --watch

2w:
	@tailwindcss -i web/main.css -o web/_main.css --watch

3w:
	@gow run cmd/main.go

# nodemon --watch './**/*.{go,templ}' --signal SIGTERM --exec 'make'