.PHONY: dev build

dev:
	@templ generate 
	tailwindcss -i ./static/css/input.css -o ./static/css/output.css
	air --build.cmd "go build -o bin/www main.go" --build.bin "./bin/www"

build:
	@templ generate
	tailwindcss -i ./static/css/input.css -o ./static/css/output.css
	go build -o bin/www main.go

.DEFAULT_GOAL := dev
