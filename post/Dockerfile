FROM golang:1.17.1-alpine3.14
WORKDIR /root/
COPY . .
RUN apk add --update --no-cache make curl

RUN make build

FROM alpine:3.14.0
RUN apk add --update --no-cache curl
WORKDIR /root/
COPY --from=0 /root/cmd/posts/posts .
EXPOSE 8002
CMD [ "./posts" ]
