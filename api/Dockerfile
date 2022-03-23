FROM golang:1.18 AS build
WORKDIR /go/src/github.com/matty-rose/espresso-chronicle/api
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin/app .

FROM scratch
WORKDIR /root/
COPY --from=build /go/src/github.com/matty-rose/espresso-chronicle/api/bin/app .
EXPOSE 8080
CMD ["./app"]
