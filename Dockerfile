FROM golang:1.10.3-alpine3.7 AS build

RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml /go/src/github.com/sluongng/rent-tracker/
WORKDIR /go/src/github.com/sluongng/rent-tracker

RUN dep ensure -vendor-only
COPY . /go/src/github.com/sluongng/rent-tracker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/rent-tracker

FROM scratch
COPY --from=build /bin/rent-tracker /bin/rent-tracker
ENTRYPOINT ["/bin/rent-tracker"]
