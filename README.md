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
Prosty serwer backendowy napisany w Go. Obsługuje podstawowe endpointy z obsługą błędów, wystarczające do przedstawienia działania aplikacji.

## Frontend (Vue3 + TypeScript)
Aplikacja kliencka typu SPA zbudowana z wykorzystaniem Composition API i Vite dla szybkiego budowania. Zapewnia responsywny interfejs użytkownika i przegląd repertuaru oraz rezerwacji biletów. Aplikacja wspiera PWA z użyciem Service Workerów.

## Uruchamianie lokalne
1. Przejdz do katalogu `backend/` i wykonaj polecenie `docker-compose up -d`
2. Przejdz do katalogu `backend/cmd/server` i uruchom serwer poleceniem `go run main.go`
3. Przejdz do katalogu `frontend/` i wykonaj polecenie `npm install`
4. Po instalacji uruchom projekt frontendu `npm run dev`
5. Otwórz `localhost:5173` w karcie przeglądarki

## Wdrażanie na DigitalOcean
1. Skopiuj `terraform/terraform.tfvars.example` do pliku `terraform/terraform.tfvars`
2. Wypełnij własnymi danymi z DigitalOcean
3. Wyeksportuj zmienne środowiskowe z danymi dostępowymi do S3 (DO Spaces, przechowywanie stanu terraform)
   * export AWS_ACCESS_KEY_ID="KEY_ID"
   * export AWS_SECRET_ACCESS_KEY="KEY_SECRET"
4. W katalogu `terraform/` wykonaj polecenie `terraform init`
5. Wykonaj polecenie `terraform plan -out tfplan`
6. Wykonaj polecenie `terraform apply "tfplan"`