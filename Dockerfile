FROM gcr.io/moonrhythm-containers/alpine

WORKDIR /app
ADD ./server .

ENTRYPOINT ["/app/server"]
