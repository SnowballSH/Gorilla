build:
	@go build ./gorilla

prod:
	@go build -ldflags "-s -w" ./gorilla

api:
	@go build  -o ./api ./gorilla/api
