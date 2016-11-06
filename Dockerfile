FROM golang:1.7-alpine

# Default list of env vars. Used them for the sake of simplcity
# Usually they should be injected in the container by another service
ENV MYSQL_HOST mysql
ENV MYSQL_USER root
ENV MYSQL_PASSWORD root
ENV MYSQL_DATABASE shorty

COPY ./ $GOPATH/src/impraise.com/shorty
RUN go install impraise.com/shorty
WORKDIR $GOPATH/src/impraise.com/shorty

CMD ["shorty"]
