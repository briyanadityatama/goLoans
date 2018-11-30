FROM golang AS build-env
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine
COPY --from=build-env /app/goLoans /usr/bin

EXPOSE 8080
ENTRYPOINT ["goLoans"]