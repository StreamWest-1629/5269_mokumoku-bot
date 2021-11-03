FROM golang:1.17

RUN apt-get update
RUN apt-get install -y ffmpeg
RUN go get github.com/bwmarrin/dca/

CMD ["/src/resource-enc/convert-in-docker.sh"]
