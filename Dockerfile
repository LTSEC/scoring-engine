FROM alpine:edge

WORKDIR ~
RUN apk --no-cache add --upgrade go curl unzip

RUN curl https://github.com/LTSEC/scoring-engine/archive/refs/heads/main.zip

RUN unzip "./main.zip"

RUN go build .

CMD [ "./scoring-engine" ]

