FROM golang
WORKDIR /LogParsing
COPY . /LogParsing
RUN go build -o LogParsing ./cmd/
COPY config/ ./config/
EXPOSE 8080
ENTRYPOINT ["./LogParsing"]