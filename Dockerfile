
#build stage
FROM golang:1.10.3-alpine
WORKDIR /go/src/app
# Copy the local package files to the container's workspace.
ADD . /go/src/gowork
COPY . .
# RUN apk add --no-cache git
# RUN go-wrapper download   # "go get -d -v ./..."
# RUN go-wrapper install    # "go install -v ./..."

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]

#final stage
# FROM alpine:1.10.3
# RUN apk --no-cache add ca-certificates
# COPY --from=builder /go/bin/app /app
# ENTRYPOINT ./app
# LABEL Name=gowork Version=0.0.1
EXPOSE 8080
