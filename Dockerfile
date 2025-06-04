# Start by building the application.
FROM golang:1.24.3-alpine as build

WORKDIR /go/src/app
COPY . .

RUN ls
RUN cat go.mod && CGO_ENABLED=0 go build -o /go/bin/app

# Now copy it into our base image.
FROM alpine

COPY --from=build /go/bin/app /

RUN apk add curl python3 ffmpeg
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /bin/yt-dlp
RUN chmod a+rx /bin/yt-dlp  # Make executable

CMD ["/app", "--cookiefile", "$COOKIES"]