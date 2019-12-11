FROM golang:latest 
RUN mkdir /app
ADD go.mod /app
ADD go.sum /app
WORKDIR /app
RUN go mod download
ADD . /app
RUN go build -o main cmd/jobsearchserver/main.go

CMD ["/app/main"]
EXPOSE 1323