FROM golang:1.8-alpine

RUN apk add --no-cache git

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download
RUN go-wrapper install -ldflags "-X main.gitCommit=`git rev-parse HEAD` -X main.buildTimestamp=`date -Iseconds`"

CMD ["go-wrapper", "run"]
