FROM golang:1.22.2-alpine3.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /simplematchmaking cmd/simplematchmaking/main.go 

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./internal/...

# Deploy the application binary into a lean image
FROM alpine:3.19 AS build-release-stage

WORKDIR /

COPY --from=build-stage /simplematchmaking /simplematchmaking

EXPOSE 8081

#USER nonroot:nonroot

ENTRYPOINT ["/simplematchmaking"]
