FROM golang:1.16.2

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
WORKDIR /workspace
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build -o .build/line-webhook-pubsub -ldflags "-w -s" .

FROM gcr.io/moonrhythm-containers/go-scratch

WORKDIR /app
COPY --from=0 /workspace/.build/* ./

ENTRYPOINT ["/app/line-webhook-pubsub"]
