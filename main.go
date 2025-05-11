package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type HistoricalFigure struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Period      string `json:"period"`
	Innovation  string `json:"innovation"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	Rank        int    `json:"rank"`
	UniqueID    string `json:"uniqueId" gorm:"uniqueIndex"`
}

type UserRanking struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FigureID  uint      `json:"figureId"`
	Rank      int       `json:"rank"`
	Timestamp time.Time `json:"timestamp"`
}

type AverageRanking struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FigureID  uint      `json:"figureId"`
	Name      string    `json:"name"`
	Average   float64   `json:"average"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func main() {
	// Initialize database
	db, err := gorm.Open(sqlite.Open("euro.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&HistoricalFigure{}, &UserRanking{}, &AverageRanking{})

	// Seed the database with historical figures
	SeedDatabase()

	// Initialize average rankings
	updateAverageRankings(db)

	// Initialize Gin router
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// API routes
	api := r.Group("/api")
	{
		api.GET("/figures", func(c *gin.Context) {
			var figures []HistoricalFigure
			db.Order("rank").Find(&figures)
			c.JSON(http.StatusOK, figures)
		})

		// Submit user ranking
		api.POST("/rankings", func(c *gin.Context) {
			var rankings []UserRanking
			if err := c.BindJSON(&rankings); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Add timestamp to each ranking
			now := time.Now()
			for i := range rankings {
				rankings[i].Timestamp = now
			}

			// Save rankings
			for _, ranking := range rankings {
				if err := db.Create(&ranking).Error; err != nil {
					log.Printf("Error creating ranking: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				log.Printf("New ranking added - Figure ID: %d, Rank: %d, Time: %s",
					ranking.FigureID, ranking.Rank, ranking.Timestamp.Format(time.RFC3339))
			}

			// Update average rankings
			updateAverageRankings(db)

			c.JSON(http.StatusOK, gin.H{"message": "Rankings submitted successfully"})
		})

		// Get average rankings
		api.GET("/rankings/average", func(c *gin.Context) {
			var averages []AverageRanking
			result := db.Order("average ASC").Find(&averages)

			// If no averages found, recalculate them
			if result.RowsAffected == 0 {
				updateAverageRankings(db)
				db.Order("average ASC").Find(&averages)
			}

			c.JSON(http.StatusOK, averages)
		})
	}

	// Start server
	fmt.Println("Starting server on port 8090")
	r.Run(":8090")
}

func updateAverageRankings(db *gorm.DB) {
	// Calculate new averages
	var averages []struct {
		FigureID uint
		Name     string
		Average  float64
	}

	db.Raw(`
		SELECT 
			ur.figure_id as figure_id,
			hf.name as name,
			AVG(ur.rank) as average
		FROM user_rankings ur
		JOIN historical_figures hf ON ur.figure_id = hf.id
		GROUP BY ur.figure_id, hf.name
	`).Scan(&averages)

	// Clear existing average rankings
	db.Exec("DELETE FROM average_rankings")

	// Create new average rankings
	for _, avg := range averages {
		newAvg := AverageRanking{
			FigureID:  avg.FigureID,
			Name:      avg.Name,
			Average:   avg.Average,
			UpdatedAt: time.Now(),
		}
		if err := db.Create(&newAvg).Error; err != nil {
			log.Printf("Error creating average ranking: %v", err)
		} else {
			log.Printf("New average ranking created - Figure: %s, Average: %.2f",
				avg.Name, avg.Average)
		}
	}

	// If no rankings exist, create default rankings based on historical figures
	if len(averages) == 0 {
		var figures []HistoricalFigure
		db.Order("rank").Find(&figures)

		for _, figure := range figures {
			newAvg := AverageRanking{
				FigureID:  figure.ID,
				Name:      figure.Name,
				Average:   float64(figure.Rank),
				UpdatedAt: time.Now(),
			}
			if err := db.Create(&newAvg).Error; err != nil {
				log.Printf("Error creating default average ranking: %v", err)
			} else {
				log.Printf("Default average ranking created - Figure: %s, Average: %.2f",
					figure.Name, float64(figure.Rank))
			}
		}
	}
}
