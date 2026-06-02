run:
	@templ generate
	@bash -c 'set -a; [ -f .env ] && source .env; set +a; go run cmd/main.go'

build:
	@templ generate 
	@go build -o ./tmp/main cmd/main.go

tailwind:
	@npx @tailwindcss/cli -i ./static/styles.css -o ./static/tailwind.css --minify