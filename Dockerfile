FROM docker.io/library/golang:1.20.2 as build
WORKDIR /

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go mod download

RUN go build -o ./bin/kubeception-agent ./cmd/agent/main.go && \
    chmod +x ./bin/kubeception-agent

RUN go build -o ./bin/kubeception-server ./cmd/server/main.go && \
    chmod +x ./bin/kubeception-server

FROM scratch

ENV PATH "$PATH:/bin"

COPY --from=build ./bin/kubeception-agent /bin/kubeception-agent
COPY --from=build ./bin/kubeception-server /bin/kubeception-server

EXPOSE 443
EXPOSE 1080