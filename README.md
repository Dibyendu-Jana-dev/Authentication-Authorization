# Technology Uses:-
  *Golang
  mongodb
  redis
  #RestApi
  ZAP for log
  #jwt Authorization
# Architecture
  Hexagonal

# Run Server:
 go run cmd/main.go

Here we implement Swagger for auto documentation generate
 Note: if you want to generate documentation all first you need to go through with the command for swagger init
 swag init -g cmd/main.go
then
 go run cmd/main.go
# Run Swagger UI
http://localhost:8080/swagger/index.html
Note: Here we implement logger with Zap package
#prerequisites
for running this service you need to install mongodb and redis must 12
