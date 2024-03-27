FROM alpine:latest

RUN "apk update && apk add go git"

RUN "git clone https://github.com/LTSEC/scoring-engine/release/tag/latest"
RUN "go build ."

CMD [ "./main" ]

