# builder image
FROM golang:1.18.4-alpine3.16 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o api .


# generate clean, final image for end users
FROM alpine:latest
COPY --from=builder /build/api .

EXPOSE 8080

# executable
ENTRYPOINT [ "./api" ]