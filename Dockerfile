FROM node:12.16.1-alpine3.11 AS js_builder
COPY webapp /webapp
WORKDIR /webapp
RUN npm i && npm run build

FROM golang:1.13.9-alpine AS go_builder
# enable modules
ENV GO111MODULE=on

WORKDIR /server
COPY server /server
# download modules
RUN go mod download

# build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build /server/cmd/server/main.go

FROM scratch 
# dont copy directly to / because static file handling would also serve linux dirs
COPY --from=go_builder /server/main /app/main
COPY --from=js_builder /webapp/build* /app/
ENV PORT=80 PRODUCTION=
EXPOSE 80
ENTRYPOINT ["/app/main"]
