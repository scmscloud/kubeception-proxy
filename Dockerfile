FROM docker.io/library/golang:1.20.2 as build
WORKDIR /

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go mod download
RUN go build -o ./bin/kubeception-proxy ./main.go

RUN chmod +x ./bin/kubeception-proxy

FROM scratch
WORKDIR /

COPY --from=build /bin/kubeception-proxy ./

EXPOSE 443
ENTRYPOINT [ "/kubeception-proxy" ]