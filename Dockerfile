FROM golang:1.17-buster

WORKDIR /app

# COPY ./go.mod .
# COPY ./go.sum .
COPY . .

RUN go mod tidy

RUN go install github.com/cosmtrek/air@v1.27.3
CMD [ "air" ]
