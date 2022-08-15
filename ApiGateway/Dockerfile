FROM golang:1.18rc1-alpine3.15
RUN mkdir api
COPY . /api
WORKDIR /api
#RUN go mod init
RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ./main