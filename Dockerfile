FROM golang:1.17

WORKDIR /go/
RUN apt-get update \
    && apt-get install -y git golang busybox

COPY . /project/

RUN go version \
    && echo ">> Checking project directories" \
    && ls -l /project/ \
    && cd /project/ \
    && echo ">> Download dependencies" \
    && go get -t -v ./... \
    && CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /project/cmd/app ./cmd \
    && echo ">> Checking builded binary" \
    && ls -lh /project/cmd/app

FROM alpine:latest

COPY --from=0 /project/cmd/app /bin/app

RUN ls -l /bin/app \
    && chmod a+rx /bin/app

USER nobody
CMD ["/bin/app"]

#add supervisor
#https://eminetto.medium.com/monitoring-a-golang-application-with-supervisor-da09cc18b498