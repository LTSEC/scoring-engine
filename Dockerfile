FROM alpine:edge

WORKDIR ~
RUN apk --no-cache add --upgrade go curl unzip

RUN curl -LJO https://github.com/LTSEC/scoring-engine/archive/refs/heads/main.zip

RUN unzip ./scoring-engine-main.zip

WORKDIR ./scoring-engine-main
RUN go build ./cmd/main.go
RUN ls -l

CMD [ "./scoring-engine-main/scoring-engine" ]

