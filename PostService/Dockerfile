FROM golang:1.18rc1-alpine3.15
RUN mkdir post
COPY . /post
WORKDIR /post
#RUN go mod init github.com/venomuz/crm-go
RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ./main