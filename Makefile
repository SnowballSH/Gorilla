build:
	@go build ./gorilla

prod:
	@go build -ldflags "-s -w" ./gorilla
