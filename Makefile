run:
	@templ generate &
	@tailwindcss -i tailwind.css -o web/public/main.css &
	@go run cmd/main.go

watch:
	@templ generate --watch &
	@tailwindcss -i tailwind.css -o web/public/main.css --watch &
	@gow run cmd/main.go

1:
	@templ generate --watch

2:
	@tailwindcss -i tailwind.css -o web/public/main.css --watch

3:
	@gow run cmd/main.go

# nodemon --watch './**/*.{go,templ}' --signal SIGTERM --exec 'make'