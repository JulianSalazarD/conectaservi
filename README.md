# ConectaServi — módulos `catalog` y `web`

Módulos del proyecto **ConectaServi** entregados en las evidencias SENA **GA7-220501096-AA2-EV01** (API REST JSON) y **GA7-220501096-AA2-EV02** (interfaz web HTML).

- **Stack**: Go 1.22 · PostgreSQL 15 · `database/sql` + driver `pgx` · Gin · `html/template`
- **Arquitectura**: monolito modular con *package by feature* — features `catalog` (dominio + CRUD JSON) y `web` (UI HTML server-rendered)
- **Estilo**: *accept interfaces, return structs* (Pike) — handlers reciben `XxxRepository` (interface), `NewPgXxxRepo` devuelve struct concreto. La feature `web` reutiliza esos repositorios sin duplicar lógica de dominio.

## Estructura

```
conectaservi/
├── cmd/api/main.go                    bootstrap (godotenv + sql.DB + gin + catalog.Mount)
├── internal/catalog/
│   ├── category.go / service.go / portfolio.go     entidades + validaciones
│   ├── errors.go                                   sentinel errors del feature
│   ├── repository.go                               interfaces CategoryRepository / ServiceRepository / PortfolioRepository
│   ├── category_repo.go / service_repo.go / portfolio_repo.go     impl Postgres (no-exportadas)
│   ├── dto.go                                      Create/Update*Request con tags `binding`
│   ├── response.go                                 writeError: mapea sentinels → 404 / 409 / 422 / 500
│   ├── category_handler.go / service_handler.go / portfolio_handler.go
│   ├── module.go                                   New(db) + Mount(r) — wiring repos→handlers→rutas
│   └── *_test.go                                   tests unitarios de entidad / VO
├── internal/web/                     EV02 — feature web (HTML + formularios server-side)
│   ├── module.go                                   New(db) + Mount(r) — registra rutas y /static
│   ├── home.go                                     GET /  (landing)
│   ├── category_pages.go / service_pages.go       handlers GET/POST por entidad
│   ├── render.go                                   render(c,name,data) + userMessage(err)
│   ├── templates.go                                //go:embed de plantillas y assets
│   ├── templates/*.html                            html/template (equivalente JSP)
│   ├── static/styles.css                           CSS plano servido en /static
│   └── category_pages_test.go                     httptest + repo fake (6 tests)
├── pkg/database/postgres.go          Open(ctx, dsn) *sql.DB con pool
├── db/schema.sql                     DDL + seed (1 user + 1 provider)
└── .env.example
```

## Quickstart

```bash
# 1. Postgres local
docker run -d --name pg \
  -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=conectaservi \
  -p 5432:5432 postgres:15
sleep 3
psql -h localhost -U postgres -d conectaservi -f db/schema.sql

# 2. Configurar conexión
cp .env.example .env
# editar DATABASE_URL si hace falta

# 3. Correr el servidor (API + UI web)
go run ./cmd/api
# [GIN-debug] Listening and serving HTTP on :8080
# Abrir http://localhost:8080/  →  panel HTML
```

## Interfaz web (EV02)

Páginas server-rendered con `html/template` (equivalente Go de JSP). Formularios HTML envían datos al servidor con métodos GET y POST. Reutiliza los repositorios del feature `catalog`.

| Método | Ruta                              | Acción                            |
|--------|-----------------------------------|-----------------------------------|
| GET    | `/`                               | Home                              |
| GET    | `/web/categories`                 | Lista de categorías               |
| GET    | `/web/categories/new`             | Formulario de creación            |
| POST   | `/web/categories`                 | Procesa creación                  |
| GET    | `/web/categories/:id/edit`        | Formulario de edición             |
| POST   | `/web/categories/:id`             | Procesa edición                   |
| POST   | `/web/categories/:id/delete`      | Eliminación                       |
| GET    | `/web/services`                   | Lista de servicios                |
| GET    | `/web/services/new`               | Formulario de creación            |
| POST   | `/web/services`                   | Procesa creación                  |
| POST   | `/web/services/:id/delete`        | Eliminación                       |
| GET    | `/static/styles.css`              | Asset estático (CSS)              |

Validación: si el dominio rechaza la entrada (slug inválido, precio negativo…), el handler **re-renderiza el mismo formulario** con los valores ingresados y el mensaje de error en español. Tras un éxito redirige (303) a la lista con un mensaje *flash*.

## Endpoints REST (EV01) — 14

| Método | Ruta                                        | CRUD                   |
|--------|---------------------------------------------|------------------------|
| POST   | `/api/v1/categories`                        | Crear                  |
| GET    | `/api/v1/categories`                        | Listar                 |
| GET    | `/api/v1/categories/:id`                    | Leer                   |
| PUT    | `/api/v1/categories/:id`                    | Actualizar             |
| DELETE | `/api/v1/categories/:id`                    | Eliminar               |
| POST   | `/api/v1/services`                          | Crear                  |
| GET    | `/api/v1/services?category_id=&is_active=`  | Listar + filtros       |
| GET    | `/api/v1/services/:id`                      | Leer                   |
| PUT    | `/api/v1/services/:id`                      | Actualizar             |
| DELETE | `/api/v1/services/:id`                      | Eliminar               |
| POST   | `/api/v1/services/:id/portfolio`            | Crear                  |
| GET    | `/api/v1/services/:id/portfolio`            | Listar                 |
| PUT    | `/api/v1/services/:id/portfolio/:pid`       | Actualizar             |
| DELETE | `/api/v1/services/:id/portfolio/:pid`       | Eliminar               |

Mapeo de errores:

| Error                         | HTTP |
|-------------------------------|------|
| `Err*NotFound`                | 404  |
| `ErrDuplicateSlug`, `ErrCategoryHasServices` | 409 |
| `ErrInvalid*`, binding fail   | 422  |
| Otros                         | 500  |

## Demo curl

```bash
B=localhost:8080/api/v1

# Crear category
curl -X POST $B/categories -H 'Content-Type: application/json' \
     -d '{"nombre":"Plomería","slug":"plomeria"}'
# → 201 {"ID":"...","Nombre":"Plomería","Slug":"plomeria",...}

# Slug duplicado
curl -X POST $B/categories -H 'Content-Type: application/json' \
     -d '{"nombre":"Plomería","slug":"plomeria"}'
# → 409 {"error":"category slug already exists"}

# Crear service (usa el provider seed y el CID de arriba)
curl -X POST $B/services -H 'Content-Type: application/json' \
     -d '{"provider_id":"00000000-0000-0000-0000-000000000001",
          "category_id":"<CID>","titulo":"Plomería 24h",
          "precio_base":50000,"lat":4.65,"lng":-74.08,"radio_km":10}'

# Intentar borrar category con services asociados → RESTRICT
curl -X DELETE $B/categories/<CID>
# → 409 {"error":"category has associated services"}

# Crear portfolio anidado
curl -X POST $B/services/<SID>/portfolio -H 'Content-Type: application/json' \
     -d '{"storage_url":"https://example.com/antes.jpg","titulo":"Antes","orden":1}'

# Borrar service → CASCADE borra sus portfolio_items
curl -X DELETE $B/services/<SID>
curl $B/services/<SID>/portfolio     # → []
```

## Tests

```bash
go test ./...
# ok  github.com/JulianSalazarD/conectaservi/internal/catalog
# ok  github.com/JulianSalazarD/conectaservi/internal/web   (httptest sobre repo fake)
```

## Estándares aplicados

| Elemento       | Convención                         | Ejemplo                                |
|----------------|------------------------------------|----------------------------------------|
| Variables      | `camelCase`, acrónimos en mayúscula| `categoryID`, `storageURL`, `createdAt`|
| Métodos        | `PascalCase` exportados            | `CategoryHandler.Create`               |
| Structs        | `PascalCase` singular              | `Category`, `Service`, `Module`        |
| Interfaces     | `PascalCase`, rol                  | `CategoryRepository`                   |
| Paquetes       | minúsculas, una palabra            | `catalog`, `database`                  |
| Archivos       | `snake_case.go`                    | `category_handler.go`                  |
| Errores        | prefijo `Err`                      | `ErrCategoryNotFound`                  |
