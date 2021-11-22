FROM golang:1.17

RUN apt-get update
RUN apt-get install -y ffmpeg
RUN go install github.com/bwmarrin/dca/cmd/dca@latest

CMD ["/src/resource-enc/convert-in-docker.sh"]
