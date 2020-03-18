FROM golang:1.13 as build
WORKDIR /build
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o duckyapi ./cmd/duckyapi/main.go


FROM alpine:latest
# Maintainer
LABEL maintainer="Ron Compos <composr@netapp.com>"
# Copy alternate fortunes
COPY --from=build /build/cmd/duckyapi/duckyapi /usr/bin/
COPY repos.txt /etc/
# Expose port 8080
EXPOSE 8080
# Run
ENTRYPOINT ["/usr/bin/duckyapi"]
