FROM golang:1.12 as builder

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -mod=vendor -o oidc-app .

FROM alpine:latest  

RUN apk --no-cache add ca-certificates && apk --no-cache add curl

WORKDIR /app/

COPY --from=builder /app/oidc-app .
COPY --from=builder /app/index.html .
COPY --from=builder /app/static/ ./static/

CMD ["/app/oidc-app"]