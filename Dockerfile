FROM golang AS builder

ENV GO111MODULE=on

RUN mkdir -p /go/src/TypeaheadBackend
WORKDIR /go/src/TypeaheadBackend

ADD . /go/src/TypeaheadBackend

WORKDIR /go/src/TypeaheadBackend

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -s -extldflags "-static"' .

FROM alpine

RUN mkdir /app 
WORKDIR /app
COPY --from=builder /go/src/TypeaheadBackend/RealDevTypeaheadBackend .

COPY names.json .

CMD ["/app/RealDevTypeaheadBackend"]
