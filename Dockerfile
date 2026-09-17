# Build stage
FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app .

# Runtime stage
FROM alpine:3.19
WORKDIR /code
COPY --from=build /bin/app /bin/app
EXPOSE 5000
CMD ["/bin/app"]
