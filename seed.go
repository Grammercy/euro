package main

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func generateUniqueID(name string, period string) string {
	// Convert to lowercase and remove spaces
	name = strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	period = strings.ToLower(strings.ReplaceAll(period, " ", "_"))
	return fmt.Sprintf("%s_%s", name, period)
}

func SeedDatabase() {
	db, err := gorm.Open(sqlite.Open("euro.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&HistoricalFigure{})

	// Clear existing data
	db.Exec("DELETE FROM historical_figures")

	figures := []HistoricalFigure{
		{
			Name:        "Johannes Gutenberg",
			Period:      "1398-1468",
			Innovation:  "Movable Type Printing Press",
			Description: "Revolutionized the spread of knowledge by inventing the movable type printing press, making books more accessible and affordable.",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/3/33/Gutenberg.jpg/800px-Gutenberg.jpg",
			Rank:        1,
			UniqueID:    generateUniqueID("Johannes Gutenberg", "1398-1468"),
		},
		{
			Name:        "Leonardo da Vinci",
			Period:      "1452-1519",
			Innovation:  "Renaissance Art and Engineering",
			Description: "Pioneered new techniques in art and engineering, creating masterpieces like the Mona Lisa and designing innovative machines.",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/b/ba/Leonardo_self.jpg",
			Rank:        2,
			UniqueID:    generateUniqueID("Leonardo da Vinci", "1452-1519"),
		},
		{
			Name:        "Nicolaus Copernicus",
			Period:      "1473-1543",
			Innovation:  "Heliocentric Theory",
			Description: "Challenged the geocentric model of the universe by proposing the heliocentric theory, revolutionizing astronomy.",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/f/f2/Nikolaus_Kopernikus.jpg",
			Rank:        3,
			UniqueID:    generateUniqueID("Nicolaus Copernicus", "1473-1543"),
		},
		{
			Name:        "Galileo Galilei",
			Period:      "1564-1642",
			Innovation:  "Scientific Method and Astronomy",
			Description: "Pioneered the scientific method and made groundbreaking discoveries in astronomy using the telescope.",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/d/d4/Justus_Sustermans_-_Portrait_of_Galileo_Galilei%2C_1636.jpg",
			Rank:        4,
			UniqueID:    generateUniqueID("Galileo Galilei", "1564-1642"),
		},
		{
			Name:        "Isaac Newton",
			Period:      "1643-1727",
			Innovation:  "Laws of Motion and Calculus",
			Description: "Formulated the laws of motion and universal gravitation, and co-invented calculus.",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/3/39/GodfreyKneller-IsaacNewton-1689.jpg",
			Rank:        5,
			UniqueID:    generateUniqueID("Isaac Newton", "1643-1727"),
		},
		{
			Name:        "James Watt",
			Period:      "1736-1819",
			Innovation:  "Steam Engine Improvements",
			Description: "Improved the steam engine, making it more efficient and practical, which was crucial for the Industrial Revolution.",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/2/2c/James_Watt_by_Henry_Howard.jpg/800px-James_Watt_by_Henry_Howard.jpg",
			Rank:        6,
			UniqueID:    generateUniqueID("James Watt", "1736-1819"),
		},
		{
			Name:        "Louis Pasteur",
			Period:      "1822-1895",
			Innovation:  "Germ Theory and Vaccination",
			Description: "Developed the germ theory of disease and created the first vaccines, revolutionizing medicine.",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e3/Louis_Pasteur_%281822_-_1895%29%2C_microbiologist_and_chemist_Wellcome_V0026980.jpg/800px-Louis_Pasteur_%281822_-_1895%29%2C_microbiologist_and_chemist_Wellcome_V0026980.jpg",
			Rank:        7,
			UniqueID:    generateUniqueID("Louis Pasteur", "1822-1895"),
		},
		{
			Name:        "Marie Curie",
			Period:      "1867-1934",
			Innovation:  "Radioactivity Research",
			Description: "Pioneered research in radioactivity, discovering polonium and radium, and won two Nobel Prizes.",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/7/7e/Marie_Curie_c1920.jpg",
			Rank:        8,
			UniqueID:    generateUniqueID("Marie Curie", "1867-1934"),
		},
		{
			Name:        "Albert Einstein",
			Period:      "1879-1955",
			Innovation:  "Theory of Relativity",
			Description: "Developed the theory of relativity, revolutionizing our understanding of space, time, and gravity.",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/3/3e/Einstein_1921_by_F_Schmutzer_-_restoration.jpg",
			Rank:        9,
			UniqueID:    generateUniqueID("Albert Einstein", "1879-1955"),
		},
		{
			Name:        "Werner Heisenberg",
			Period:      "1901-1976",
			Innovation:  "Quantum Mechanics",
			Description: "Formulated the uncertainty principle and made fundamental contributions to quantum mechanics.",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/f/f8/Bundesarchiv_Bild183-R57262%2C_Werner_Heisenberg.jpg",
			Rank:        10,
			UniqueID:    generateUniqueID("Werner Heisenberg", "1901-1976"),
		},
	}

	// Use FirstOrCreate to prevent duplicates
	for _, figure := range figures {
		var existing HistoricalFigure
		result := db.Where("unique_id = ?", figure.UniqueID).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&figure).Error; err != nil {
				log.Printf("Error creating figure: %v", err)
			} else {
				log.Printf("Created new figure: %s", figure.Name)
			}
		} else {
			log.Printf("Figure already exists: %s", figure.Name)
		}
	}

	log.Println("Database seeded successfully!")
}
