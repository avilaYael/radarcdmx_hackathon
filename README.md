<div align="center">

# 🛰️ Radar CDMX - https://radarcdmx.maykelfarha.com/
# 👥 Equipo - Axodata

### Inteligencia territorial y económica para la Ciudad de México

Dashboard de análisis territorial con **indicadores en vivo** y **detección automática de inconsistencias de uso de suelo**, sobre un solo mapa que unifica datos públicos hoy dispersos.

<br>

`Hackathon SecretarIA` · `SEDECO + Saptiva AI` · **Reto 1 — Análisis territorial de datos**

<br>

![Stack](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![DB](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql&logoColor=white)
![Map](https://img.shields.io/badge/MapLibre_GL-Mapas-396CB2)
![Datos](https://img.shields.io/badge/Datos-100%25_p%C3%BAblicos-9F2241)
![Estado](https://img.shields.io/badge/estado-funcional_end--to--end-2e7d32)


Backend generado usando [nuzur](https://nuzur.com/)

<br>

<!-- ╔══════════════════════════════════════════════════════════════════╗
     ║  👇 CAMBIAR ESTOS 3 LINKS ANTES DE ENVIAR EL FORM (16:00–17:00)  ║
     ╚══════════════════════════════════════════════════════════════════╝ -->

[**🌐 Demo en vivo**](https://radarcdmx.maykelfarha.com/) ·
[**⚙️ Guía de instalación →**](SETUP.md)

</div>

---

## 1. Problema en una frase

La información geográfica y económica de la Ciudad de México —establecimientos del DENUE, mercados públicos y uso de suelo SEDUVI— vive dispersa en fuentes, bases de datos y archivos físicos distintos, lo que obliga a SEDECO a invertir días en localizar, cruzar y validar datos antes de poder decidir una política económica.

## 2. A quién sirve

> **Usuario primario:** el/la **analista de la Dirección de Política Económica de SEDECO** que necesita responder, en minutos y sin solicitar extracciones a otra área, preguntas como:
> *"¿Cuántos restaurantes hay en Cuauhtémoc, cuántos empleos representan y cuáles operan en un uso de suelo incompatible?"*

**Usuarios secundarios:** personal de ventanilla que valida la compatibilidad de uso de suelo de un establecimiento, y equipos de planeación que comparan alcaldías para priorizar programas.

### Capacidades clave

| | Capacidad | Qué resuelve |
|---|-----------|--------------|
| 🗺️ | **Mapa unificado** | DENUE + mercados públicos + zonificación SEDUVI en una sola vista interactiva |
| 🔎 | **Búsqueda geoespacial** | Establecimientos por radio, alcaldía, sector (SCIAN), actividad y uso de suelo |
| ⚖️ | **Dictamen de uso de suelo** | Veredicto automático *compatible / incompatible* del giro contra el uso de suelo, con su razón |
| 📊 | **Comparador de alcaldías** | Total de establecimientos, empleos estimados y mezcla por sector económico |

## 3. Cómo correrlo

> Todo el producto —API + interfaz— lo sirve **un solo binario de Go** en `http://localhost:4242`.
> **Requisitos:** [Go 1.26+](https://go.dev/dl/) y Docker.

```bash
# 1) Clonar
git clone https://github.com/IsraelGonzalezCruz/Hacton-cdmxV1.git   # 👈 CAMBIAR si el repo se renombra
cd Hacton-cdmxV1

# 2) Base de datos + esquema
docker run --name radarcdmx-mysql -e MYSQL_ROOT_PASSWORD=radar -e MYSQL_DATABASE=radarcdmx \
  -p 3306:3306 -d mysql:8
docker exec -i radarcdmx-mysql mysql -uroot -pradar radarcdmx \
  < backend/rcapi/core/repository/sql/schema/create.sql

# 3) Configurar conexión  (plantilla lista en backend/dev.example.yaml)
cp backend/dev.example.yaml backend/dev.yaml

# 4) Cargar datos de demo + arrancar
cd frontend
CONFIG=../backend/dev.yaml go run ./cmd/import -csv ../data/establecimientos_seed.csv
CONFIG=../backend/dev.yaml ADDR=:4242 go run .
```

Abre **http://localhost:4242** y haz clic en cualquier punto del mapa para ver su dictamen de uso de suelo.

📖 **Pasos detallados, resolución de problemas y estructura del proyecto → [`SETUP.md`](SETUP.md).**

## 4. Stack usado

| Capa | Tecnología |
|------|-----------|
| **API de dominio** (`backend/rcapi`) | Go 1.26 · gRPC + HTTP · [Uber fx](https://github.com/uber-go/fx) (DI) · [sqlc](https://sqlc.dev) · JWT |
| **Servidor web** (`frontend`) | Go 1.26 · [chi](https://github.com/go-chi/chi) · fx · embebe `rcapi` in-process y expone REST en `/api/*` |
| **Base de datos** | MySQL 8 (InnoDB) · consultas geoespaciales `ST_Distance_Sphere` · columnas `JSON` para ubicación/contacto |
| **Interfaz** (`frontend/RadarMX-main`) | HTML5 + JavaScript (sin framework) · MapLibre GL JS · Chart.js · Lucide · servida estática desde el binario Go (sin paso de build) |
| **Datos públicos** | DENUE (INEGI) · Mercados Públicos CDMX · Uso de suelo SEDUVI / SIG CDMX |
| **Despliegue** | Dockerfile multi-stage · charts Helm · CI en `.github/workflows` |

<details>
<summary><b>Arquitectura (diagrama)</b></summary>

```
Navegador  ·  MapLibre GL + dashboard
     │  HTTP /api/*
     ▼
radarcdmx-web  (Go · :4242)  ──in-process──►  rcapi core  ──►  MySQL 8
   ├─ /api/establecimientos          lista paginada
   ├─ /api/establecimientos/nearby   radio geoespacial (ST_Distance_Sphere) + filtros
   ├─ /api/dictamen-uso-de-suelo     ¿el uso de suelo es compatible con el giro?
   ├─ /api/municipios/compare        totales · empleos estimados · mezcla por sector
   ├─ /api/{sectores|actividades|municipios|usos-de-suelo|mercados|zoning}
   └─ /healthz
```

**Motor de dictamen:** del `codigo_actividad` (SCIAN) se deriva el sector a 2 dígitos y se contrasta el `uso_de_suelo` del establecimiento contra el catálogo de usos permitidos (`sectores.json`), devolviendo `aprobado` + la razón. En el cliente, además, la coordenada se cruza contra los polígonos de zonificación SEDUVI (`zoning.json`) mediante un algoritmo punto-en-polígono.

</details>

## 5. Limitaciones conocidas

- **El dictamen de uso de suelo es determinista** (reglas SCIAN ↔ usos permitidos), no usa un LLM: es correcto y auditable, pero no interpreta lenguaje natural ni los casos límite del RETYS.
- **Datos de demostración acotados:** la rentabilidad trimestral, la población por alcaldía y la tasa de cumplimiento son estimaciones/sintéticos plausibles cuando no hay endpoint dedicado; los conteos de establecimientos, los empleos por rango `per_ocu` y el dictamen **sí** salen de datos reales del DENUE cargados en MySQL.
- **Requiere MySQL 8** (por `ST_Distance_Sphere` y las funciones `JSON_*`); no corre en MySQL 5.7 ni SQLite.
- **El CSV completo del DENUE y `backend/dev.yaml` están en `.gitignore`** (peso y credenciales); la demo se reproduce con el seed de `data/` y la plantilla `backend/dev.example.yaml`.
- **La capa SEDUVI (`zoning.json`) es una muestra** de polígonos para la demo, no la cobertura catastral completa.
- **Mapbox es opcional:** sin token se usa la capa base gratuita de CartoDB; el token solo mejora el estilo del mapa.

## Datos

La demo se ejecuta 100 % sobre **información pública** (válido y esperado por las bases del reto):

| Fuente | Uso en el producto | Origen |
|--------|--------------------|--------|
| Unidades económicas DENUE | Puntos en el mapa · empleos por `per_ocu` · dictamen | INEGI |
| Mercados Públicos CDMX | Capa de mercados (`mercados.json`) | Gobierno CDMX |
| Uso de suelo SEDUVI | Polígonos (`zoning.json`) y catálogo de usos permitidos (`sectores.json`) | SIG CDMX / SEDUVI |

---

<div align="center">

Construido durante el **Hackathon SecretarIA** · SEDECO + Saptiva AI · 6 de junio de 2026
<br>
Datos públicos de INEGI, SEDUVI / SIG CDMX y Gobierno de la Ciudad de México

</div>
