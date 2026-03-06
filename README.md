# CinemaReservation

Demo projekt mikroserwisu Go, wdrażanego na DigitalOcean za pomocą Terraform i Docker.

## Infrastruktura (Terraform)
Infrastruktura jest zarządzana jako kod (IaC) i składa się z następujących plików:
* `provider.tf` - Konfiguracja providera DigitalOcean oraz backendu S3 (Spaces) do przechowywania stanu.
* `vpc.tf` - Definicja prywatnej sieci (VPC).
* `database.tf` - Zarządzalny klaster bazy danych PostgreSQL umieszczony w sieci VPC.
* `registry.tf` - Prywatny rejestr kontenerów (DOCR).
* `app.tf` - Definicja usługi w App Platform z bezpiecznym przekazaniem adresu bazy danych przez sieć wewnętrzną (`private_uri`).
* `variable.tf` - Deklaracje zmiennych wejściowych (token API, tag obrazu).

## Konteneryzacja (Dockerfile)
Aplikacja jest pakowana przy użyciu *multi-stage build* zapewniając mały rozmiar obrazu:
* **Etap budowania:** Kompilacja kodu Go w kontenerze `golang`.
* **Etap uruchomieniowy:** Wykorzystanie obrazu `distroless/base-debian11`. Aplikacja uruchamiana jest z uprawnieniami `nonroot`, bez dostępu do powłoki systemowej, co zapewnia bezpieczeństwo.

## Mikroserwis (Go)
Prosty serwer backendowy napisany w Go. 
Obecnie aplikacja posiada jeden endpoint HTTP (`/`) nasłuchujący na porcie `8080`, który służy jako health check i zwraca krótki komunikat `Running` na ekranie.
