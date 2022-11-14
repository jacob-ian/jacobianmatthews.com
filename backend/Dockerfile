# MUST BE BUILT WITH ROOT CONTEXT
# DEVELOPMENT ENVIRONMENT
FROM golang:1.19-alpine as dev
RUN apk add --update make postgresql-client

RUN go install github.com/cosmtrek/air@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /src/jacobianmatthews.com
COPY go.work . 

WORKDIR /src/jacobianmatthews.com/backend
COPY backend .

CMD ["air"]

# TESTING ENVIRONMENT
FROM golang:1.19-alpine as test
WORKDIR /src/jacobianmatthews.com/backend
COPY backend .
RUN go test ./..

# BUILD ENVIRONMENT
FROM golang:1.19-alpine as build
WORKDIR /src/jacobianmatthews.com
COPY go.work . 
WORKDIR /src/jacobianmatthews.com/backend
COPY backend .
RUN go build -o bin/backend cmd/backend/main.go 

# PRODUCTION ENVIRONMENT
FROM alpine:3.16 as prod
WORKDIR /opt/jacobianmatthews.com/backend
COPY --from=build /src/jacobianmatthews.com/backend/bin bin
CMD ["./bin/backend"]
