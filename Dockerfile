FROM golang:1.20 as go
FROM gcr.io/distroless/base-debian10 as run

FROM go as build
WORKDIR /go/src/github.com/sakajunquality/clouddeploy-functions

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /go/bin/app server/main.go server/controller.go

FROM run
COPY --from=build /go/bin/app /usr/local/bin/app
CMD ["app"]
