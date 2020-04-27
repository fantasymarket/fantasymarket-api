FROM golang:alpine AS build
RUN mkdir -p /var/FantasyMarket
ADD . /var/FantasyMarket
WORKDIR /var/FantasyMarket
RUN go build -ldflags="-s -w" -i main.go -o bin/fantasymarket


FROM alpine
ARG VERSION=V1.0
ARG PORT=3000
LABEL com.FantasyMarket.version=$VERSION
ENV GOLANG_ENV="production"
COPY --from=build /var/FantasyMarket /var/FantasyMarket
WORKDIR /var/FantasyMarket
EXPOSE $PORT
ENTRYPOINT [".bin/fantasymarket"]

