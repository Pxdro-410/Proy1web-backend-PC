# F1 Pilot Tracker — Backend

## Autor: Pedro Caso 241286
API REST construida con **Go** que expone los datos de pilotos de Fórmula 1 y su sistema de ratings. No genera HTML: devuelve únicamente JSON. Cualquier cliente (navegador, app móvil, script) puede consumirla de la misma manera, deployada en Railway.

> **Repositorio del frontend:** https://github.com/Pxdro-410/Proy1web-frontend-PC

---

## Screenshot del backend deployado y funcionando

<img width="1600" height="830" alt="image" src="https://github.com/user-attachments/assets/49944710-cf6f-41c7-9ad0-ad929e90bebf" />


---

## Stack tecnológico

| Capa | Tecnología |
|---|---|
| Lenguaje | Go 1.22+ |
| Base de datos | PostgreSQL (Railway) |
| Driver de BD | `github.com/lib/pq` |
| Deploy | Railway |

---

## Cómo correr el proyecto localmente

### Pre-requisitos

- [Go 1.22+](https://go.dev/dl/) instalado
- Una instancia de PostgreSQL corriendo (local o en la nube)

### 1. Clonar el repositorio

```bash
git clone https://github.com/Pxdro-410/Proy1web-backend-PC.git
cd Proy1web-backend-PC
```

### 2. Configurar la base de datos

**Opción A — Con `DATABASE_URL` como Railway:**

```powershell
# PowerShell
$env:DATABASE_URL = "postgresql://usuario:contraseña@host:puerto/nombre_db"
```

**Opción B — Con variables individuales en PostgreSQL local:**

```powershell
$env:DB_HOST     = "localhost"
$env:DB_PORT     = "5432"
$env:DB_USER     = "postgres"
$env:DB_PASSWORD = "postgres"
$env:DB_NAME     = "f1tracker"
```

> Si no defines ninguna variable, el servidor usará los valores por defecto de la Opción B.

### 3. Instalar dependencias y ejecutar

```bash
go mod download
go run .
```

El servidor arranca en `http://localhost:8080`.

---

## CORS

El frontend y el backend corren en dominios distintos. El navegador bloquea las peticiones fetch() a menos que el servidor lo permita explícitamente mediante headers CORS. CORS lo que hace controlar las peticiones de una pagina web hacia hacia un dominio distinto al suyo.

El backend está configurado para aceptar peticiones de cualquier origen:

**Configuración aplicada:**

```
Access-Control-Allow-Origin:  *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type
```

Se permite cualquier origen durante desarrollo. En un entorno de producción real se restringiría al dominio del cliente.

---

## Endpoints de la API

### Pilotos

| Método | Ruta | Descripción | Código éxito |
|---|---|---|---|
| `GET` | `/piloto` | Listar pilotos (paginado) | 200 |
| `GET` | `/piloto/:id` | Obtener piloto por ID | 200 |
| `POST` | `/piloto` | Crear piloto nuevo | 201 |
| `PUT` | `/piloto/:id` | Editar piloto existente | 200 |
| `DELETE` | `/piloto/:id` | Eliminar piloto | 204 |

### Query params de `GET /piloto`

| Parámetro | Tipo | Descripción | Ejemplo |
|---|---|---|---|
| `page` | int | Número de página (default: 1) | `?page=2` |
| `limit` | int | Resultados por página, máx. 100 (default: 10) | `?limit=6` |
| `q` | string | Búsqueda por nombre (case-insensitive) | `?q=hamilton` |
| `sort` | string | Campo a ordenar: `name`, `team`, `number`, `championships`, `created_at` | `?sort=championships` |
| `order` | string | Dirección: `ASC` o `DESC` | `?order=DESC` |

### Ratings

| Método | Ruta | Descripción | Código éxito |
|---|---|---|---|
| `GET` | `/piloto/:id/rating` | Obtener ratings y promedio del piloto | 200 |
| `POST` | `/piloto/:id/rating` | Enviar un rating al piloto | 201 |

### Respuesta de error en formato JSON

```json
{
  "error": "descripción del error"
}
```

### Respuesta paginada (`GET /piloto`)

```json
{
  "data": [...],
  "page": 1,
  "limit": 12,
  "total": 15,
  "total_pages": 2
}
```

---

## Challenges implementados

| Challenge | Puntos | Descripción |
|---|---|---|
| **Códigos HTTP correctos** | 20 pts | `201` al crear, `204` al eliminar, `404` si no existe, `400` en input inválido |
| **Validación server-side** | 20 pts | Campos requeridos, rangos numéricos, respuestas de error en JSON descriptivo (`{ "error": "..." }`) |
| **Paginación** | 30 pts | `GET /piloto?page=&limit=` con metadatos: `total`, `total_pages`, `page`, `limit` |
| **Búsqueda** | 15 pts | `?q=` con `ILIKE` en PostgreSQL, case-insensitive, compatible con paginación |
| **Ordenamiento** | 15 pts | `?sort=` con whitelist de columnas + `?order=asc\|desc`, validado server-side |
| **Sistema de rating** | 30 pts | Tabla `ratings` propia en BD, endpoints REST dedicados, promedio calculado en el servidor |
| **Total** | **130 pts** | |

---

## Reflexión

### Go como lenguaje de backend

Se eligió go por la eficiencia y poca complegidad para la creacion de las APIs y conexión a la DB. Permite construir un servidor HTTP production-ready sin ningún framework externo. El sistema de tipos estático ayudó a detectar errores en tiempo de compilación antes de que llegaran a producción. Y por ultimo pero no menos importante la velocidad de compilación y ejecución es notablemente superior a otros lenguajes más complejos.

**¿Lo usaría de nuevo?** Sí, definitivamente para proyectos donde el rendimiento y la simplicidad del deploy importen.

### PostgreSQL como base de datos

PostgreSQL fue más que suficiente para este proyecto. Es una tecnología excelente para el manejo de la información de la base de datos, en este caso, no se exprime el máximo rendimiento de Postgres debido a la simplicidad de la estructura de las entidades.

**¿Lo usaría de nuevo?** Si, ya lo he usado anteriormente y es excelente para el manejo de base de datos relacionales y además aprendí que en Railway se provee de forma gratuita para el deploy.

### Railway como plataforma de deploy

Aprendí que Railway detecta automáticamente los proyectos mediante el `go.mod` y los despliega sin configuración adicional. Es la plataforma más sencilla que he usado para desplegar un backend de forma gratuita.

**¿Lo usaría de nuevo?** Si, me pareció una plataforma bastante intuitiva que definitivamente usaré para otros proyectos de backend.
