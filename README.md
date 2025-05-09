# Top 10 Innovators in AP European History

An interactive web application showcasing the most innovative figures in European history. This project features a modern, responsive design with interactive elements and a timeline view.

## Features

- Interactive timeline of historical figures
- Detailed profiles with images and descriptions
- Responsive design for all devices
- Modern UI with smooth animations
- RESTful API backend

## Prerequisites

- Go 1.21 or higher
- Modern web browser

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd euro
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run main.go
```

4. Open your browser and navigate to:
```
http://localhost:8080
```

## Project Structure

- `main.go` - Main application entry point
- `seed.go` - Database seeding script
- `templates/` - HTML templates
- `static/` - Static assets (CSS, JavaScript, images)
  - `css/` - Stylesheets
  - `js/` - JavaScript files

## Technologies Used

- Backend:
  - Go (Gin framework)
  - GORM (ORM)
  - SQLite (Database)

- Frontend:
  - HTML5
  - CSS3 (with Tailwind CSS)
  - JavaScript (Vanilla)
  - Responsive Design

## Contributing

Feel free to submit issues and enhancement requests! 