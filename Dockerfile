# BUILD
FROM golang:latest AS builder
ADD . /app
WORKDIR /app
RUN go mod download
RUN ls -la
WORKDIR /app/cmd/webservice
# CGO and link flags are required in this case for sqllite
# Other flags (-w -s) being used to minimise resultant image size
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -linkmode external -extldflags -static" -a -o ./main .

# RUN
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/cmd/webservice/main /app/main
COPY --from=builder /app/cmd/webservice/assets /app/assets
COPY --from=builder /app/cmd/webservice/web /app/web
WORKDIR /app
RUN chmod +x ./main
ENTRYPOINT [ "./main" ]
EXPOSE 8080