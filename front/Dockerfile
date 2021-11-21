FROM denoland/deno:alpine
WORKDIR /root/
COPY . .
RUN apk add --update --no-cache make curl

EXPOSE 8080
# Using development mode 'cause of DOMParser error
RUN make dev
