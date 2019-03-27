FROM golang:1.11-alpine

WORKDIR /app
COPY . .

# TODO do multystage build and run bin app
CMD go run app.go