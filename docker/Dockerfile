FROM golang:1.22 AS build

WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /nisit-backend ./cmd/app/main.go

FROM gcr.io/distroless/base-debian11 AS release

WORKDIR /

COPY --from=build /nisit-backend /nisit-backend

# USER user:user

ENTRYPOINT [ "/nisit-backend" ]
