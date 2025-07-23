# URL Shortener (Bit.ly Clone)

A modern, full-stack URL shortener application built with **Golang** backend and **Next.js** frontend. Create short links from long URLs, track click statistics, and analyze link performance.

![URL Shortener](https://img.shields.io/badge/Status-Ready-green.svg)
![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)
![Next.js](https://img.shields.io/badge/Next.js-14+-black.svg)
![TypeScript](https://img.shields.io/badge/TypeScript-5+-blue.svg)

## ✨ Features

- **🔗 URL Shortening**: Convert long URLs into short, memorable links
- **📊 Analytics**: Track clicks, monitor link performance with detailed statistics
- **🎨 Custom Codes**: Create custom short codes for branded links
- **⚡ Fast Redirects**: Redis caching for lightning-fast redirects
- **📱 Responsive Design**: Beautiful, modern UI that works on all devices
- **🔍 Click Tracking**: Monitor clicks by date, country, and referrer
- **🚀 High Performance**: Built with Go for optimal speed and efficiency

## 🛠️ Tech Stack

### Backend
- **Golang** - High-performance backend language
- **Gin** - Fast HTTP web framework
- **GORM** - Object-relational mapping
- **PostgreSQL** - Primary database
- **Redis** - Caching layer for fast redirects
- **Base62 Encoding** - Short code generation

### Frontend
- **Next.js 14** - React framework with SSR
- **TypeScript** - Type-safe JavaScript
- **Tailwind CSS** - Utility-first CSS framework
- **Axios** - HTTP client for API communication
- **React Hot Toast** - Beautiful notifications
- **Lucide React** - Modern icon library
- **Chart.js** - Data visualization (for statistics)

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+**
- **Node.js 18+**
- **PostgreSQL 13+**
- **Redis 6+** (optional, for caching)

### Option 1: Docker (Recommended)

1. **Clone the repository**
```bash
git clone <repository-url>
cd shorter
```

2. **Start with Docker Compose**
```bash
docker-compose up -d
```

3. **Access the application**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Health Check: http://localhost:8080/health

### Option 2: Manual Setup

#### Backend Setup

1. **Navigate to backend directory**
```bash
cd backend
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. **Set up PostgreSQL database**
```sql
CREATE DATABASE shorter_db;
```

5. **Run the backend**
```bash
go run main.go
```

The backend will start on `http://localhost:8080`

#### Frontend Setup

1. **Navigate to frontend directory**
```bash
cd frontend
```

2. **Install dependencies**
```bash
npm install
```

3. **Set up environment variables**
```bash
cp .env.local.example .env.local
# Edit .env.local if needed
```

4. **Run the frontend**
```bash
npm run dev
```

The frontend will start on `http://localhost:3000`

## 📋 Environment Variables

### Backend (.env)
```env
# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=shorter_db
DB_PORT=5432
DB_SSLMODE=disable

# Redis Configuration (Optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Server Configuration
PORT=8080
GIN_MODE=debug
CUSTOM_DOMAIN=localhost:8080
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

## 🔧 API Endpoints

### URLs
- `POST /api/shorten` - Create a short URL
- `GET /api/urls` - Get all URLs (paginated)
- `GET /api/stats/:code` - Get click statistics for a URL
- `GET /:code` - Redirect to original URL

### Health
- `GET /health` - Health check endpoint

## 📖 Usage

### Creating Short URLs

1. Visit the homepage
2. Enter your long URL in the input field
3. Optionally add a custom code
4. Click "Shorten URL"
5. Copy and share your short link!

### Viewing Statistics

1. Find your URL in the "Recent URLs" section
2. Click the statistics button
3. View detailed analytics including:
   - Total clicks
   - Daily click trends
   - Geographic distribution
   - Referrer sources

## 🏗️ Project Structure

```
shorter/
├── backend/                 # Go backend application
│   ├── config/             # Database and Redis configuration
│   ├── handlers/           # HTTP request handlers
│   ├── models/             # Database models
│   ├── utils/              # Utility functions
│   ├── main.go             # Application entry point
│   ├── go.mod              # Go dependencies
│   └── Dockerfile          # Backend Docker configuration
├── frontend/               # Next.js frontend application
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── pages/          # Next.js pages
│   │   ├── styles/         # CSS styles
│   │   ├── types/          # TypeScript type definitions
│   │   └── utils/          # Utility functions
│   ├── package.json        # Node.js dependencies
│   └── Dockerfile          # Frontend Docker configuration
├── docker-compose.yml      # Docker Compose configuration
└── README.md              # Project documentation
```

## 🧪 Development

### Backend Development
```bash
cd backend
go run main.go
```

### Frontend Development
```bash
cd frontend
npm run dev
```

### Building for Production

#### Backend
```bash
cd backend
go build -o main .
```

#### Frontend
```bash
cd frontend
npm run build
npm start
```

## 🐳 Docker Deployment

Build and run with Docker Compose:

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild services
docker-compose up -d --build
```

## 🔒 Security Features

- Input validation and sanitization
- SQL injection prevention with GORM
- Rate limiting (can be added)
- CORS configuration
- Safe URL validation

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [Next.js](https://nextjs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [GORM](https://gorm.io/)
- [Lucide Icons](https://lucide.dev/)

## 📞 Support

If you have any questions or need help, please:

1. Check the [Issues](../../issues) page
2. Create a new issue if needed
3. Contact the maintainers

---

**Built with ❤️ using Go and Next.js** # shorter
