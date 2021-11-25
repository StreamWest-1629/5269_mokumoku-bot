FROM golang:1.17

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .
RUN apt-get update

# encoder
RUN apt-get install -y ffmpeg
RUN go install github.com/bwmarrin/dca/cmd/dca@latest
# debugger
# RUN go install github.com/cosmtrek/air@v1.27.3

COPY . .
RUN go mod tidy
RUN go build .


CMD [ /app/app ]
