FROM golang:1.19.2
RUN mkdir /wb_l0 
ADD . /wb_l0/ 
WORKDIR /wb_l0 
RUN go build /wb_l0/cmd/app/main.go 
CMD ["/app/main"]