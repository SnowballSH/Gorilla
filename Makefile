build:
	@go build .

prod:
	@go build -ldflags "-s -w" .
