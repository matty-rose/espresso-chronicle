FROM golang:1.18 AS build
WORKDIR /go/src/github.com/matty-rose/espresso-chronicle/api
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify
COPY app ./app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin/app app/main.go

FROM golang:1.18 as dev
ENV GO111MODULE=on
WORKDIR /go/src/github.com/matty-rose/espresso-chronicle/api
COPY --from=build /go/src/github.com/matty-rose/espresso-chronicle/api/ ./
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
CMD CompileDaemon --build="go build app/main.go" --command "./main" --color

FROM scratch as prod
WORKDIR /root/
COPY --from=build /go/src/github.com/matty-rose/espresso-chronicle/api/bin/app .
EXPOSE 8000
CMD ["./app"]
