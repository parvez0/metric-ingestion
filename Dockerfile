FROM golang
MAINTAINER syedparvez72@gmail.com
# creating directories for application files and db data directory
RUN mkdir -p /{app, data}
WORKDIR /app
COPY . /app
# environment variables with their default values
ENV PORT 8080
ENV LOGLEVEL info
ENV SQLITE_DB_PATH=/data
ENV GO_ENV=production
# volume mount path
VOLUME /data
# testing image
RUN go test -v test_scripts
# building go executable
RUN go build -o app
# creating users and making it owner of the required directories
RUN groupadd go \
    && useradd gouser -G go \
    && chown -R gouser:go /app \
    && chown -R gouser:go /data \
    && chmod -x app

# starting the application
CMD ["./app"]