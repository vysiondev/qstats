FROM golang:1.15-alpine AS build
RUN apk update && apk add --no-cache git
RUN apk --no-cache add ca-certificates

WORKDIR /build/qstats

COPY go.mod ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/qstats

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/qstats /bin/qstats
COPY --from=build /build/qstats/conf /bin/conf

WORKDIR /bin
EXPOSE 3030
ENTRYPOINT ["qstats"]