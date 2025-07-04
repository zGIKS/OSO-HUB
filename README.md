# OSOHUB Backend (Go + Cassandra + Swagger)

Este proyecto es un backend seguro y profesional en Go para OSOHUB, usando Cassandra como base de datos y documentación Swagger generada con swag CLI.

---

## Requisitos previos

```bash
sudo apt update
sudo apt install golang-go
```

---

## Instalación y setup

1. **Clona el repositorio y entra a la carpeta backend:**

```bash
cd /home/mateo/Documentos/imgur/backend
```

2. **Inicializa el módulo de Go:**

```bash
go mod init osohub
```

3. **Instala dependencias:**

```bash
go get github.com/gocql/gocql
# Framework web seguro y rápido
go get github.com/gin-gonic/gin
# Swagger para documentación
go get github.com/swaggo/gin-swagger
# Archivos estáticos de Swagger
go get github.com/swaggo/files
# Carga variables de entorno de .env (opcional, recomendado)
go get github.com/joho/godotenv
```

4. **Instala swag CLI (para generar docs Swagger):**

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

5. **Crea el archivo `.env` en `backend/` con tus variables:**

```
CASSANDRA_HOST=127.0.0.1
CASSANDRA_KEYSPACE=osohub
PORT=8080
```

6. **Genera la documentación Swagger:**

```bash
swag init --generalInfo cmd/main.go --output docs
```

7. **Ejecuta el backend:**

```bash
go run cmd/main.go
```

Swagger estará disponible en: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## Estructura recomendada

```
backend/
├── cmd/
│   └── main.go
├── config/
│   └── config.go
├── db/
│   └── cassandra.go
├── handlers/
│   └── users.go
│   └── images.go
│   └── ...
├── models/
│   └── user.go
│   └── image.go
│   └── ...
├── docs/         # Swagger docs generados
├── go.mod
├── go.sum
└── README.md
```

---

## Buenas prácticas y seguridad

- Uso de variables de entorno para credenciales y configuración.
- Manejo de errores seguro y mensajes claros.
- Validación de datos de entrada (UUID, emails, etc).
- Uso de context y timeouts en queries a Cassandra.
- Documentación Swagger para todos los endpoints.
- Separación de lógica en capas (handlers, models, db).
- No exponer información sensible en logs ni respuestas.

---

## Ejemplo de conexión segura a Cassandra (`db/cassandra.go`)

```go
package db

import (
	"log"
	"os"
	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

var Session *gocql.Session

func InitCassandra() {
	_ = godotenv.Load() // Carga .env si existe
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))
	cluster.Keyspace = os.Getenv("CASSANDRA_KEYSPACE")
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second
	cluster.ConnectTimeout = 5 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Error connecting to Cassandra: %v", err)
	}
	Session = session
}
```

---

## Ejemplo de handler seguro (`handlers/users.go`)

```go
package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"osohub/db"
	"osohub/models"
)

// GetUserByID godoc
// @Summary Get user by ID
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /users/{user_id} [get]
func GetUserByID(c *gin.Context) {
	idStr := c.Param("user_id")
	userID, err := gocql.ParseUUID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var user models.User
	query := `SELECT user_id, username, email, password_hash, profile_picture_url, bio, role, created_at FROM users_by_id WHERE user_id = ? LIMIT 1`
	if err := db.Session.Query(query, userID).Consistency(gocql.One).Scan(
		&user.UserID, &user.Username, &user.Email, &user.PasswordHash,
		&user.ProfilePictureURL, &user.Bio, &user.Role, &user.CreatedAt,
	); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
```

---

## Seguridad extra
- Usa HTTPS en producción.
- No loguees contraseñas ni hashes.
- Limita el tamaño de los requests.
- Usa middlewares para autenticación/autorización.

---

## Para más endpoints y modelos
Sigue el patrón mostrado en `handlers/` y `models/` para las otras tablas (`images_by_date`, `images_by_user`, etc).

---

¡Listo! Tu backend Go seguro, profesional y documentado con Swagger para Cassandra está preparado.
