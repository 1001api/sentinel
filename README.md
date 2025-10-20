![banner](https://i.imgur.com/UjWZ1rh.png)

# Sentinel

A powerful, self-hosted event tracking and analytics platform built with Go and Fiber. Sentinel provides real-time insights into user behavior, event tracking, and comprehensive analytics for your web applications.

## ✨ Features

- **Real-time Event Tracking** - Capture and analyze user interactions in real-time
- **Geolocation Intelligence** - Automatic IP-based geolocation using MaxMind GeoIP2
- **Session Management** - Track user sessions with Redis-backed storage
- **API Key Management** - Secure API access with key-based authentication
- **Data Export** - Export events to Excel and PDF formats
- **Live Analytics Dashboard** - Beautiful web interface built with Templ, Alpine.js, and TailwindCSS
- **Rate Limiting** - Built-in protection against abuse
- **Worker Pool** - Efficient background processing for high-volume events
- **OAuth Integration** - Google OAuth2 authentication support
- **Caching Layer** - Redis-based caching for optimal performance
- **RESTful API** - Public and private API endpoints for event ingestion and retrieval

## 🏗️ Architecture

Sentinel follows a clean architecture pattern with:

- **Fiber** - High-performance web framework
- **PostgreSQL** - Primary data store with pgx driver
- **Redis** - Session management and caching
- **Templ** - Type-safe Go templating engine
- **SQLC** - Type-safe SQL query generation
- **Worker Pool** - Concurrent event processing
- **GeoIP2** - IP geolocation database

## 📋 Prerequisites

- Go 1.23.2 or higher
- PostgreSQL 12+
- Redis 6+
- Node.js 18+ (for frontend assets)
- MaxMind GeoLite2 database (included but not auto-updated)

## 🚀 Getting Started

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/hubkudev/sentinel.git
   cd sentinel
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Install Node.js dependencies**
   ```bash
   npm install
   ```

4. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. **Run database migrations**
   ```bash
   # Migrations are automatically applied on startup
   ```

6. **Download GeoIP2 database**
   - Place `GeoLite2-City.mmdb` in `internal/ipdb/` directory
   - Sign up at [MaxMind](https://www.maxmind.com/) for free database access

### Development

Run the development server with hot-reload:

```bash
make live
```

This starts 5 parallel processes:
- Templ generation with live reload
- Go server with Air hot-reload
- TailwindCSS compilation
- TypeScript/JavaScript bundling with esbuild
- Asset synchronization

Or run individual processes:

```bash
make live/server    # Go server only
make live/templ     # Templ generation only
make live/tailwind  # CSS compilation only
make live/esbuild   # JS bundling only
```

### Production Build

```bash
# Build frontend assets
npm run build
npx tailwindcss -i ./views/public/input.css -o ./views/public/global.css --minify

# Generate Templ templates
templ generate

# Build Go binary
go build -o sentinel .

# Run
./sentinel
```

## 📡 API Usage

### Public Event Tracking

Send events to Sentinel:

```bash
POST /api/v1/event
Content-Type: application/json
X-API-Key: your-public-api-key

{
  "ProjectID": "uuid-here",
  "EventType": "page_view",
  "EventLabel": "Homepage",
  "PageURL": "https://example.com",
  "UserAgent": "Mozilla/5.0...",
  "BrowserName": "Chrome",
  "SessionID": "session-id",
  "DeviceType": "desktop",
  "TimeOnPage": 5000,
  "ScreenResolution": "1920x1080",
  "FiredAt": "2024-10-20T11:22:00Z"
}
```

### Retrieve Events

```bash
GET /api/v1/events?project_id=uuid-here
X-API-Key: your-private-api-key
```

## 🗂️ Project Structure

```
sentinel/
├── configs/          # Configuration files (DB, Redis, GeoIP)
├── gen/             # SQLC generated code
├── internal/
│   ├── constants/   # Application constants and validators
│   ├── dto/         # Data transfer objects
│   ├── entities/    # Domain entities
│   ├── ipdb/        # GeoIP2 database files
│   ├── middlewares/ # HTTP middlewares
│   ├── migrations/  # Database migrations
│   ├── mocks/       # Test mocks
│   ├── repositories/# Data access layer
│   ├── routes/      # HTTP route definitions
│   ├── schema/      # SQL schemas for SQLC
│   └── services/    # Business logic layer
├── views/
│   ├── components/  # Reusable Templ components
│   ├── pages/       # Page templates
│   ├── public/      # Static assets (CSS, JS)
│   └── static/      # TypeScript source files
├── main.go          # Application entry point
├── Makefile         # Development commands
└── sqlc.yaml        # SQLC configuration
```

## 🔧 Configuration

Key environment variables:

```env
PORT=8080
DATABASE_URL=postgresql://user:pass@localhost:5432/sentinel
REDIS_URL=redis://localhost:6379
COOKIE_SALT=your-secret-key
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test file
go test ./internal/services/user_service_test.go
```

## 📊 Features in Detail

### Event Tracking
- Capture custom events with flexible schema
- Automatic IP geolocation (country, region, city)
- Session tracking and user identification
- Browser and device detection
- Page URL and element path tracking

### Analytics Dashboard
- Real-time event monitoring
- Weekly event charts
- Event type and label distribution
- Geographic visitor distribution
- Browser and device analytics
- Session analysis

### Data Export
- Export events to Excel (XLSX)
- Generate PDF reports
- Customizable date ranges
- Filtered exports by project

### API Management
- Create multiple API keys per project
- Public keys for event ingestion
- Private keys for data retrieval
- Key-based rate limiting

## 🛡️ Security

- Encrypted cookies using AES
- API key authentication
- Rate limiting on sensitive endpoints
- CORS protection
- Session-based authentication for web interface
- OAuth2 integration

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Web framework
- [Templ](https://templ.guide/) - Go templating
- [SQLC](https://sqlc.dev/) - SQL code generation
- [MaxMind](https://www.maxmind.com/) - GeoIP2 database
- [TailwindCSS](https://tailwindcss.com/) - CSS framework
- [Alpine.js](https://alpinejs.dev/) - JavaScript framework

## 📧 Support

For issues and questions, please open an issue on GitHub.

---

Built with ❤️ using Go and Fiber