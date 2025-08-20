package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"

	_ "modernc.org/sqlite" // Driver SQLite sin CGO
)

// -------------------- MODELOS --------------------
type Doll struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Nombre string `json:"nombre"`
	Edad   int    `json:"edad"`
	Activo bool   `json:"activo"`
	Cartas int    `json:"cartas"`
}

type Cliente struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Nombre   string `json:"nombre"`
	Ciudad   string `json:"ciudad"`
	Motivo   string `json:"motivo"`
	Contacto string `json:"contacto"`
}

type Carta struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	ClienteID uint   `json:"cliente_id"`
	DollID    uint   `json:"doll_id"`
	Fecha     string `json:"fecha"`
	Estado    string `json:"estado"` // borrador, revisado, enviado
	Contenido string `json:"contenido"`
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
	db.AutoMigrate(&Doll{}, &Cliente{}, &Carta{})
	log.Println("✅ Conectado a SQLite (modernc.org/sqlite)")
}

// -------------------- MAIN --------------------
func main() {
	initDB()
	r := gin.Default()

	// -------------------- CRUD DOLLS --------------------
	r.GET("/dolls", func(c *gin.Context) {
		var dolls []Doll
		db.Find(&dolls)
		c.JSON(http.StatusOK, dolls)
	})

	r.POST("/dolls", func(c *gin.Context) {
		var doll Doll
		if err := c.ShouldBindJSON(&doll); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
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
		if err := c.ShouldBindJSON(&doll); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
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

	// -------------------- CRUD CLIENTES --------------------
	r.GET("/clientes", func(c *gin.Context) {
		var clientes []Cliente
		db.Find(&clientes)
		c.JSON(http.StatusOK, clientes)
	})

	r.POST("/clientes", func(c *gin.Context) {
		var cliente Cliente
		if err := c.ShouldBindJSON(&cliente); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		if err := c.ShouldBindJSON(&cliente); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
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
		db.Find(&cartas)
		c.JSON(http.StatusOK, cartas)
	})

	r.POST("/cartas", func(c *gin.Context) {
		var carta Carta
		if err := c.ShouldBindJSON(&carta); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&carta)
		c.JSON(http.StatusOK, carta)
	})

	r.PUT("/cartas/:id", func(c *gin.Context) {
		var carta Carta
		id := c.Param("id")
		if err := db.First(&carta, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carta no encontrada"})
			return
		}
		if err := c.ShouldBindJSON(&carta); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&carta)
		c.JSON(http.StatusOK, carta)
	})

	r.DELETE("/cartas/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&Carta{}, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carta no encontrada"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"mensaje": "Carta eliminada"})
	})

	// -------------------- SERVIDOR --------------------
	r.Run(":8080")
}
