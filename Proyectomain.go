package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// -------------------- MODELOS --------------------
type Doll struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Nombre string `json:"nombre" gorm:"not null"`
	Edad   int    `json:"edad" gorm:"not null"`
	Activo bool   `json:"activo"`
	Cartas int    `json:"cartas"` // contador total histórico
}

type Cliente struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Nombre   string `json:"nombre" gorm:"not null"`
	Ciudad   string `json:"ciudad" gorm:"not null"`
	Motivo   string `json:"motivo" gorm:"not null"`
	Contacto string `json:"contacto" gorm:"not null"`
}

type Carta struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	ClienteID uint   `json:"cliente_id"`
	DollID    uint   `json:"doll_id"`
	Fecha     string `json:"fecha" gorm:"not null"`
	Estado    string `json:"estado" gorm:"not null"` // borrador, revisado, enviado
	Contenido string `json:"contenido" gorm:"not null"`
}

// -------------------- VARIABLES GLOBALES --------------------
var db *gorm.DB

// -------------------- CONEXIÓN A DB --------------------
func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("violet.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Error conectando a SQLite:", err)
	}
	if err := db.AutoMigrate(&Doll{}, &Cliente{}, &Carta{}); err != nil {
		log.Fatal("❌ Error migrando modelos:", err)
	}
	log.Println("✅ Conectado a SQLite (pure Go driver)")
}

// -------------------- MIDDLEWARE CORS --------------------
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// -------------------- HELPERS --------------------
func validateEstado(estado string) bool {
	switch estado {
	case "borrador", "revisado", "enviado":
		return true
	default:
		return false
	}
}

func countActiveLettersForDoll(dollID uint) int64 {
	var count int64
	db.Model(&Carta{}).
		Where("doll_id = ? AND estado IN ?", dollID, []string{"borrador", "revisado"}).
		Count(&count)
	return count
}

func findAvailableDoll() (Doll, bool) {
	var dolls []Doll
	db.Where("activo = ?", true).Order("id ASC").Find(&dolls)
	for _, d := range dolls {
		if countActiveLettersForDoll(d.ID) < 5 {
			return d, true
		}
	}
	return Doll{}, false
}

// -------------------- MAIN --------------------
func main() {
	initDB()
	r := gin.Default()
	r.Use(cors())

	// -------------------- CRUD DOLLS --------------------
	r.GET("/dolls", func(c *gin.Context) {
		var dolls []Doll
		db.Find(&dolls)
		c.JSON(http.StatusOK, dolls)
	})

	r.POST("/dolls", func(c *gin.Context) {
		var doll Doll
		if err := c.ShouldBindJSON(&doll); err != nil || strings.TrimSpace(doll.Nombre) == "" || doll.Edad <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos para la Doll"})
			return
		}
		doll.Cartas = 0
		db.Create(&doll)
		c.JSON(http.StatusOK, doll)
	})

	r.PUT("/dolls/:id", func(c *gin.Context) {
		var doll Doll
		id := c.Param("id")
		if err := db.First(&doll, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Doll no encontrada"})
			return
		}
		var input Doll
		if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Nombre) == "" || input.Edad <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}
		doll.Nombre = input.Nombre
		doll.Edad = input.Edad
		doll.Activo = input.Activo
		db.Save(&doll)
		c.JSON(http.StatusOK, doll)
	})

	r.DELETE("/dolls/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&Doll{}, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Doll no encontrada"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"mensaje": "Doll eliminada"})
	})

	// -------------------- CRUD CLIENTES (con búsquedas) --------------------
	r.GET("/clientes", func(c *gin.Context) {
		var clientes []Cliente
		q := db.Model(&Cliente{})
		if nombre := c.Query("nombre"); nombre != "" {
			q = q.Where("LOWER(nombre) LIKE ?", "%"+strings.ToLower(nombre)+"%")
		}
		if ciudad := c.Query("ciudad"); ciudad != "" {
			q = q.Where("LOWER(ciudad) LIKE ?", "%"+strings.ToLower(ciudad)+"%")
		}
		q.Find(&clientes)
		c.JSON(http.StatusOK, clientes)
	})

	r.POST("/clientes", func(c *gin.Context) {
		var cliente Cliente
		if err := c.ShouldBindJSON(&cliente); err != nil ||
			strings.TrimSpace(cliente.Nombre) == "" ||
			strings.TrimSpace(cliente.Ciudad) == "" ||
			strings.TrimSpace(cliente.Motivo) == "" ||
			strings.TrimSpace(cliente.Contacto) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos para el cliente"})
			return
		}
		db.Create(&cliente)
		c.JSON(http.StatusOK, cliente)
	})

	r.PUT("/clientes/:id", func(c *gin.Context) {
		var cliente Cliente
		id := c.Param("id")
		if err := db.First(&cliente, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
			return
		}
		var input Cliente
		if err := c.ShouldBindJSON(&input); err != nil ||
			strings.TrimSpace(input.Nombre) == "" ||
			strings.TrimSpace(input.Ciudad) == "" ||
			strings.TrimSpace(input.Motivo) == "" ||
			strings.TrimSpace(input.Contacto) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}
		cliente.Nombre = input.Nombre
		cliente.Ciudad = input.Ciudad
		cliente.Motivo = input.Motivo
		cliente.Contacto = input.Contacto
		db.Save(&cliente)
		c.JSON(http.StatusOK, cliente)
	})

	r.DELETE("/clientes/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&Cliente{}, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"mensaje": "Cliente eliminado"})
	})

	// -------------------- CRUD CARTAS --------------------
	r.GET("/cartas", func(c *gin.Context) {
		var cartas []Carta
		q := db.Model(&Carta{})
		if cliente := c.Query("cliente"); cliente != "" {
			if cid, err := strconv.Atoi(cliente); err == nil {
				q = q.Where("cliente_id = ?", cid)
			}
		}
		if estado := c.Query("estado"); estado != "" {
			q = q.Where("estado = ?", estado)
		}
		q.Find(&cartas)
		c.JSON(http.StatusOK, cartas)
	})

	r.POST("/cartas", func(c *gin.Context) {
		var body struct {
			ClienteID uint   `json:"cliente_id"`
			Fecha     string `json:"fecha"`
			Contenido string `json:"contenido"`
		}
		if err := c.ShouldBindJSON(&body); err != nil ||
			body.ClienteID == 0 || strings.TrimSpace(body.Fecha) == "" || strings.TrimSpace(body.Contenido) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos para la carta"})
			return
		}

		// Asignación automática de Doll activa con < 5 cartas en proceso
		doll, ok := findAvailableDoll()
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No hay Dolls disponibles (máximo 5 cartas activas por Doll)"})
			return
		}

		carta := Carta{
			ClienteID: body.ClienteID,
			DollID:    doll.ID,
			Fecha:     body.Fecha,
			Estado:    "borrador",
			Contenido: body.Contenido,
		}
		if err := db.Create(&carta).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la carta"})
			return
		}
		// contador histórico
		db.Model(&Doll{}).Where("id = ?", doll.ID).UpdateColumn("cartas", gorm.Expr("cartas + 1"))
		c.JSON(http.StatusOK, carta)
	})

	r.PUT("/cartas/:id", func(c *gin.Context) {
		var carta Carta
		id := c.Param("id")
		if err := db.First(&carta, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carta no encontrada"})
			return
		}
		var input struct {
			Estado    string `json:"estado"`
			Contenido string `json:"contenido"`
		}
		if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Contenido) == "" || !validateEstado(input.Estado) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		// Validar flujo de estados: borrador -> revisado -> enviado
		valid := (carta.Estado == "borrador" && input.Estado == "revisado") ||
			(carta.Estado == "revisado" && input.Estado == "enviado")
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Transición de estado inválida (borrador→revisado→enviado)"})
			return
		}

		carta.Estado = input.Estado
		carta.Contenido = input.Contenido
		db.Save(&carta)
		c.JSON(http.StatusOK, carta)
	})

	r.DELETE("/cartas/:id", func(c *gin.Context) {
		var carta Carta
		id := c.Param("id")
		if err := db.First(&carta, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carta no encontrada"})
			return
		}
		if carta.Estado != "borrador" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden eliminar cartas en estado borrador"})
			return
		}
		db.Delete(&carta)
		c.JSON(http.StatusOK, gin.H{"mensaje": "Carta eliminada"})
	})

	// -------------------- REPORTES --------------------
	r.GET("/reportes/dolls/:id", func(c *gin.Context) {
		id := c.Param("id")
		var totalCartas int64
		var clientesDistintos int64
		db.Model(&Carta{}).Where("doll_id = ?", id).Count(&totalCartas)
		db.Model(&Carta{}).Where("doll_id = ?", id).Distinct("cliente_id").Count(&clientesDistintos)
		var doll Doll
		db.First(&doll, id)
		c.JSON(http.StatusOK, gin.H{
			"doll_id":            doll.ID,
			"nombre":             doll.Nombre,
			"total_cartas":       totalCartas,
			"clientes_distintos": clientesDistintos,
		})
	})

	// -------------------- SERVIDOR --------------------
	r.Run(":8080")
}
