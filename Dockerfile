# Budowanie
FROM golang:1.25 AS build-stage

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

# Kopiowanie wszystkiego z zachowaniem struktury katalogów
COPY . .

# Budowanie z wskazaniem na pakiet main
RUN CGO_ENABLED=0 GOOS=linux go build -o /cinema-reservation ./cmd/server

# Finalny obraz (distroless)
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

# Kopiowanie binarki do obrazu runtime
COPY --from=build-stage /cinema-reservation /cinema-reservation

# Aplikacja nasłuchuje na porcie
EXPOSE 8080

# Domyślny użytkownik bez uprawnień roota (distroless
USER nonroot:nonroot

# Uruchomienie serwisu
CMD ["/cinema-reservation"]
