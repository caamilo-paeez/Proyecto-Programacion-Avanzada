# Compañía Postal CH — Trabajo de Avanzada

Backend en Go (Gin + GORM + SQLite puro Go) y frontend HTML/CSS/JS conectado por `fetch`.

## Requisitos
- Go 1.20+ (recomendado 1.22)
- No requiere CGO (usa `github.com/glebarez/sqlite`)

## Cómo ejecutar
```bash
# En la raíz del proyecto
go mod tidy
go run Proyectomain.go
```
El backend queda en `http://localhost:8080`.

Abre `index.html` en tu navegador (o sirve el directorio con un server estático).

## Endpoints clave
- `POST /cartas` **auto-asigna** una Doll activa con <5 cartas en proceso; la carta inicia en `borrador`.
- `PUT /cartas/:id` controla el flujo `borrador → revisado → enviado`.
- `DELETE /cartas/:id` solo permitido si la carta está en `borrador`.
- `GET /clientes?ciudad=...&nombre=...` búsquedas.
- `GET /reportes/dolls/:id` reporte por Doll.

## Postman
Importa `postman-collection.json` y prueba todos los flujos.
