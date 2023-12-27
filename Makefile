run:
	@templ generate &
	@tailwindcss -i tailwind.css -o public/main.css &
	@go run cmd/main.go

watch:
	@templ generate --watch &
	@tailwindcss -i tailwind.css -o public/main.css --watch &
	@gow run cmd/main.go

# @nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run cmd/main.go
# @templ generate --watch --proxy="http://localhost:8080" --cmd="runtest"