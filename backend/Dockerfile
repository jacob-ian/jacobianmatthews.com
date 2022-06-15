# MUST BE BUILT WITH ROOT CONTEXT
# DEVELOPMENT ENVIRONMENT
FROM golang:1.18-alpine as dev

WORKDIR /src/jacobianmatthews.com
COPY go.work . 

WORKDIR /src/jacobianmatthews.com/backend
COPY backend .

RUN go install github.com/cosmtrek/air@latest

CMD ["air"]
