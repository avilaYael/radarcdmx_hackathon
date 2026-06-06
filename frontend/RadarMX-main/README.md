# Radar CDMX | Dashboard de Análisis Territorial - SEDECO

**Radar CDMX** es una plataforma frontend interactiva diseñada para la Secretaría de Desarrollo Económico (SEDECO) de la Ciudad de México. Permite consultar la correspondencia del uso de suelo comercial contra los polígonos de zonificación de SEDUVI y realizar análisis territoriales.

## 🚀 Características
- **Auditoría de Uso de Suelo**: Algoritmo geométrico de punto-en-polígono (Raycasting) que valida si la ubicación y giro de un comercio (DENUE) cumple con el uso de suelo asignado (H, HC, I, E).
- **HUD de Filtros de Búsqueda (SCIAN 2023)**:
  - Búsqueda por Alcaldía, Sector Económico, Actividad (Código) y Uso de Suelo.
  - El selector de actividad funciona como subfiltro del sector seleccionado.
  - Al seleccionar filtros, el mapa se encuadra y desplaza suavemente (`fitBounds` / `flyTo`) hacia los establecimientos coincidentes.
  - Descripciones institucionales completas integradas.
- **Importación de Datos (CSV)**: Carga de archivos espaciales mediante arrastrar y soltar (Drag and Drop) con soporte de geolocalización.
- **Comparador Territorial**: Módulo comparativo de métricas de población, establecimientos, mercados y cumplimiento normativo entre alcaldías de la CDMX.
- **Gráficos Financieros**: Desempeño trimestral del establecimiento auditado usando Chart.js.
- **Identidad Visual**: Paleta de colores institucional del Gobierno de la CDMX (Guinda `#9F2241` y tonos de Beige) con el imagotipo "Ajolotito CDMX".

## 🛠️ Tecnologías Utilizadas
- **Core**: HTML5, Vanilla JavaScript.
- **Estilos**: CSS3 moderno (Glassmorphism, CSS Variables, animaciones premium).
- **Mapa**: MapLibre GL JS (capas vectoriales CartoDB / Mapbox).
- **Gráficos**: Chart.js.
- **Iconos**: Lucide Icons.
- **Entorno**: Vite (servidor de desarrollo estático ultrarrápido).

## 📁 Estructura del Proyecto
```text
├── dist/                    # Compilación de producción (HTML/CSS/JS optimizados)
├── src/
│   ├── api.js               # Servicios mock y lógica matemática de punto-en-polígono
│   ├── app.js               # Controlador e interacciones del mapa y de la interfaz
│   └── dataset.js           # Datos espaciales base (GeoJSONs de SEDUVI, DENUE y Mercados)
├── ajolotito.jpeg           # Logotipo e identidad de la aplicación
├── index.html               # Estructura del dashboard principal
├── style.css                # Sistema de diseño, temas y responsivo
├── package.json             # Dependencias del proyecto
└── README.md                # Esta guía de uso
```

## 💻 Instalación y Desarrollo Local

1. Asegúrate de tener instalado [Node.js](https://nodejs.org/).
2. Instala las dependencias de desarrollo (Vite):
   ```bash
   npm install
   ```
3. Inicia el servidor de desarrollo local:
   ```bash
   npm run dev
   ```
4. Abre en tu navegador la dirección que se muestre en terminal (usualmente `http://localhost:3000` o `http://localhost:5173`).

## 🐳 Integración de Backend (Guía para Desarrolladores Go / GORM)
Para que el equipo de backend migre las bases de datos espaciales y conecte la API real, consulta la **Guía de Integración** detallada que se encuentra en la carpeta del proyecto o en el archivo de diseño correspondiente. 

La guía especifica:
- Las estructuras de datos GORM en Go (`Establecimiento`, `ZonaSeduvi`, `Mercado`).
- Los endpoints REST requeridos para suplir los datos mock.
- Las consultas SQL espaciales utilizando **PostGIS** (`ST_Contains`, `ST_SetSRID`, `ST_MakePoint`).
