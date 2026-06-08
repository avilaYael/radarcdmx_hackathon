# ⚙️ SETUP — Radar CDMX

Guía probada para clonar y correr **Radar CDMX** de cero. Tiempo estimado: **~5 minutos** (más la descarga de la imagen de MySQL la primera vez).

Todo el producto —API REST + interfaz web— lo sirve **un único binario de Go** en `http://localhost:4242`.

---

## 1. Requisitos

| Herramienta | Versión | Para qué |
|-------------|---------|----------|
| [Go](https://go.dev/dl/) | **1.26+** | Compilar y correr el backend y el servidor web |
| [Docker](https://www.docker.com/) | cualquiera reciente | Levantar MySQL 8 sin instalarlo a mano |
| [Node.js](https://nodejs.org/) | 18+ *(opcional)* | Solo si quieres iterar la UI con Vite (hot-reload) |

> **¿Por qué MySQL 8 y no otra?** El motor de búsqueda por cercanía usa `ST_Distance_Sphere` y las funciones `JSON_*` de MySQL 8. No funciona en MySQL 5.7 ni en SQLite.

---

## 2. Clonar el repositorio

```bash
git clone https://github.com/mklfarha/radarcdmx  
cd radarcdmx
```

---

## 3. Base de datos

### 3.1 Levantar MySQL 8 con Docker

```bash
docker run --name radarcdmx-mysql \
  -e MYSQL_ROOT_PASSWORD=radar \
  -e MYSQL_DATABASE=radarcdmx \
  -p 3306:3306 -d mysql:8
```

Espera ~15 s a que inicialice. Verifica que esté arriba:

```bash
docker exec radarcdmx-mysql mysqladmin -uroot -pradar ping
# -> mysqld is alive
```

### 3.2 Crear el esquema (tablas `user` y `establecimiento`)

```bash
docker exec -i radarcdmx-mysql mysql -uroot -pradar radarcdmx \
  < backend/rcapi/core/repository/sql/schema/create.sql
```

---

## 4. Configurar la conexión

La configuración con credenciales (`backend/dev.yaml`) está en `.gitignore`. Crea la tuya a partir de la plantilla incluida:

```bash
cp backend/dev.example.yaml backend/dev.yaml
```

Contenido esperado de `backend/dev.yaml` (ya viene así en la plantilla; ajusta solo si cambiaste credenciales o puerto):

```yaml
ports:
  grpc: 6011
  http: 8080

auth:
  jwt:
    key: dev-secret

db:
  - name: radarcdmx
    host: 127.0.0.1
    port: "3306"
    user: root
    pswd: radar
    params: parseTime=true
    driver: mysql
```

---

## 5. Cargar los datos de demo

El importador transmite el CSV del DENUE y lo inserta en MySQL a través de `rcapi`. Usa el seed incluido en `data/`:

```bash
cd frontend
CONFIG=../backend/dev.yaml go run ./cmd/import -csv ../data/establecimientos_seed.csv
```

Flags útiles del importador:

| Flag | Default | Descripción |
|------|---------|-------------|
| `-csv` | `./establecimientos_para_importar.csv` | Ruta del CSV a importar |
| `-workers` | `10` | Inserciones concurrentes |
| `-delay` | `500us` | Pausa entre inserciones por worker (evita saturar la BD) |
| `-start` | `0` | Salta N filas (para reanudar una carga) |
| `-limit` | `0` | Detiene tras N filas (`0` = sin límite) |

> Para cargar el dataset **completo** del DENUE, descárgalo de INEGI, colócalo como `data/establecimientos_para_importar.csv` y pásalo con `-csv`.

---

## 6. Arrancar el producto

```bash
# desde frontend/
CONFIG=../backend/dev.yaml ADDR=:4242 go run .
```

Verás en consola las rutas montadas. Abre 👉 **http://localhost:4242**

- El mapa carga establecimientos (DENUE) y mercados públicos.
- **Clic en cualquier punto** → panel lateral con el **dictamen de uso de suelo** y los indicadores.
- Usa los filtros (alcaldía, sector SCIAN, actividad, uso de suelo) y el **comparador de alcaldías**.

### Hot-reload opcional (desarrollo)

```bash
# desde frontend/ — requiere github.com/air-verse/air
air     # usa .air.toml, que ya apunta a ../backend/dev.yaml
```

### Iterar la UI con Vite (opcional)

La interfaz se sirve estática desde el binario Go y **no necesita build**. Solo si quieres recarga en caliente del front:

```bash
cd frontend/RadarMX-main
npm install
npm run dev      # http://localhost:5173 (ajusta el proxy /api hacia :4242)
```

---

## 7. Estructura del proyecto

> **Nota sobre la convención de carpetas.** El brief sugiere `/src` y `/data`. Este proyecto separa el código en dos módulos Go independientes —`backend/` (API de dominio) y `frontend/` (servidor web + UI)— porque el `frontend` referencia al `backend` mediante un `replace` local en `go.mod`; renombrarlos rompería esa relación. Equivalencia: **`backend/` + `frontend/` = `/src`**. Los datos públicos de la demo viven en **`data/`**.

```
.
├── README.md                  # Resumen, problema, usuario, stack, limitaciones
├── SETUP.md                   # Esta guía
├── data/                      # Datos públicos de la demo
│   └── establecimientos_seed.csv
│
├── backend/                   # ── API de dominio (Go) ── [parte de /src]
│   ├── dev.example.yaml        # Plantilla de configuración (copiar a dev.yaml)
│   └── rcapi/
│       ├── main.go             # Servicio gRPC + HTTP (fx)
│       ├── core/               # Módulos establecimiento / user + acceso a datos (sqlc)
│       │   └── repository/sql/schema/create.sql   # Esquema MySQL
│       ├── entity/  enum/  idl/  auth/  config/
│
└── frontend/                  # ── Servidor web + interfaz (Go) ── [parte de /src]
    ├── main.go                 # Arranque del servidor (fx)
    ├── cmd/import/             # Importador CSV DENUE → MySQL
    ├── internal/
    │   ├── handlers/           # Endpoints REST /api/* (+ catálogos sectores/zoning/mercados.json)
    │   ├── server/  route/  config/  rcapi/
    └── RadarMX-main/           # Interfaz (mapa + dashboard), servida estática
        ├── index.html  style.css  ajolotito.jpeg
        └── src/  app.js  api.js
```

---

## 8. Endpoints (referencia rápida)

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/` | Interfaz web (dashboard) |
| GET | `/healthz` | Health check |
| GET | `/api/establecimientos?page_size=&offset=` | Lista paginada |
| GET | `/api/establecimientos/nearby?lat=&lng=&radius_m=&...` | Búsqueda por radio + filtros (actividad, uso de suelo, municipio) |
| GET | `/api/dictamen-uso-de-suelo?uuid=` | Veredicto de compatibilidad de uso de suelo |
| GET | `/api/municipios/compare?a=&b=` | Comparativo entre dos alcaldías |
| GET | `/api/sectores` · `/api/actividades` · `/api/municipios` · `/api/usos-de-suelo` | Catálogos |
| GET | `/api/mercados` · `/api/zoning` | Capas GeoJSON |

Prueba rápida:

```bash
curl "http://localhost:4242/healthz"
curl "http://localhost:4242/api/establecimientos?page_size=3"
```

---

## 9. Resolución de problemas

| Síntoma | Causa probable | Solución |
|---------|----------------|----------|
| `config file does not exist` / `config path is empty` | No se exportó `CONFIG` o no existe `backend/dev.yaml` | `cp backend/dev.example.yaml backend/dev.yaml` y exporta `CONFIG=../backend/dev.yaml` |
| `dial tcp 127.0.0.1:3306: connect: connection refused` | MySQL aún no levanta | Espera unos segundos o revisa `docker ps` / `docker logs radarcdmx-mysql` |
| `Error 1054 ... ST_Distance_Sphere` o funciones `JSON_*` | MySQL < 8 | Usa la imagen `mysql:8` indicada |
| `frontend asset not found` | Se ejecutó fuera de `frontend/` | Corre `go run .` **desde** `frontend/` (la UI se sirve desde `RadarMX-main/`) |
| El mapa carga pero sin puntos | No se importaron datos | Ejecuta el paso 5 (importador CSV) |
| `go: requires go >= 1.26` | Go desactualizado | Instala Go 1.26+ desde [go.dev/dl](https://go.dev/dl/) |

---

## 10. Despliegue (opcional)

Imagen Docker del servidor web (contexto = raíz del repo, por el `replace` hacia `backend/rcapi`):

```bash
docker build -f frontend/Dockerfile -t radarcdmx-web .
docker run -p 4242:4242 -e CONFIG=/app/dev.yaml -v $(pwd)/backend/dev.yaml:/app/dev.yaml radarcdmx-web
```

También se incluyen charts de **Helm** en `backend/rcapi/.helm` y `frontend/` y workflows de CI en `.github/workflows`.
