run:
	@templ generate
	@go run cmd/main.go

build:
	@templ generate 
	@go build -o ./tmp/main cmd/main.go