FROM golang:alpine AS build

WORKDIR /build
COPY v1 v1
COPY ["go.mod", "go.sum", "./"]
RUN CGO_ENABLED=0 go build -o server ./v1/cmd/server/main.go

FROM alpine:3

COPY --from=build /build/server /server

EXPOSE 80

CMD /server
