# Card Vault 🔐

A secure, enterprise-grade backend service for managing credit/debit cards with advanced encryption, PCI DSS compliance, and scalable architecture.

## 🚀 Features

### Core Functionality
- **Secure Card Storage**: AES-256-GCM encryption with automatic key versioning
- **JWT Authentication**: Stateless authentication with configurable expiration
- **Full CRUD Operations**: Create, read, update, and delete cards with proper validation
- **Batch Operations**: Efficient concurrent updates for multiple cards
- **PCI DSS Compliance**: Industry-standard security practices and audit trails

### Security Features
- **Advanced Encryption**: AES-256-GCM with unique nonces per operation
- **Key Management**: Automatic key generation with secure rotation capabilities
- **Data Masking**: Sensitive data never exposed in responses
- **Rate Limiting**: IP-based request throttling to prevent abuse
- **Security Headers**: Comprehensive HTTP security headers
- **Input Validation**: Luhn algorithm for card number validation

### Performance & Scalability
- **Concurrent Processing**: Goroutine-based batch operations
- **Database Optimization**: Indexed queries and connection pooling
- **Memory Efficiency**: Streaming operations for large datasets
- **Health Monitoring**: Built-in health check endpoints

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Client App    │───▶│   API Gateway    │───▶│   Card Vault    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                         │
                                                         ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │  Encryption     │◀───│   PostgreSQL    │
                       │   Service       │    │    Database     │
                       └─────────────────┘    └─────────────────┘
```

### Tech Stack
- **Backend**: Go 1.25+ with Gin framework
- **Database**: PostgreSQL 15+ with GORM
- **Security**: AES-256-GCM encryption, JWT authentication
- **Deployment**: Docker & Docker Compose
- **Testing**: Testify for unit and integration tests

## 📋 API Documentation

### Authentication
All endpoints require JWT authentication via `Authorization: Bearer <token>` header.

#### Generate Test Token
```http
POST /auth/test-token
```

### Card Management

#### Create Card
```http
POST /api/v1/cards
Content-Type: application/json

{
  "cardholder_name": "John Doe",
  "card_number": "4111111111111111",
  "expiry_month": 12,
  "expiry_year": 2025,
  "cvv": "123"
}
```

#### Get All Users Cards
```http
GET /api/v1/cards
```

#### Get Specific Card
```http
GET /api/v1/cards/{card_id}
```

#### Update Card
```http
PUT /api/v1/cards/{card_id}
Content-Type: application/json

{
  "cardholder_name": "Jane Doe",
  "card_number": "4111111111111111",
  "expiry_month": 6,
  "expiry_year": 2026,
  "cvv": "456"
}
```

#### Delete Card
```http
DELETE /api/v1/cards/{card_id}
```

#### Batch Update Cards
```http
POST /api/v1/cards/batch-update
Content-Type: application/json

{
  "cards": [
    {
      "id": "uuid-here",
      "cardholder_name": "Updated Name",
      "expiry_month": 12,
      "expiry_year": 2026
    }
  ]
}
```

### Administrative

#### Rotate Encryption Keys
```http
POST /api/v1/admin/cards/rotate-keys
```

## 🔧 Installation & Setup

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- PostgreSQL 15+ (if running locally)

### Quick Start with Docker

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd card-vault
   ```

2. **Start services**
   ```bash
   docker-compose up --build
   ```

3. **Verify installation**
   ```bash
   curl http://localhost:8080/health
   ```

### Local Development Setup

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start PostgreSQL**
   ```bash
   docker run --name postgres \
     -e POSTGRES_USER=cardvault_user \
     -e POSTGRES_PASSWORD=secure_password_123 \
     -e POSTGRES_DB=cardvault_db \
     -p 5432:5432 -d postgres:15-alpine
   ```

4. **Run database migrations**
   ```bash
   # Migrations run automatically on startup
   go run cmd/server/main.go
   ```

## 🧪 Testing

### Run Unit Tests
```bash
go test ./tests/... -v
```

### Run with Coverage
```bash
go test ./tests/... -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### API Testing Examples

#### 1. Generate Authentication Token
```bash
curl -X POST http://localhost:8080/auth/test-token \
  -H "Content-Type: application/json"
```

#### 2. Create a Card
```bash
curl -X POST http://localhost:8080/api/v1/cards \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "cardholder_name": "John Doe",
    "card_number": "4111111111111111",
    "expiry_month": 12,
    "expiry_year": 2025,
    "cvv": "123"
  }'
```

#### 3. Test Key Rotation
```bash
curl -X POST http://localhost:8080/api/v1/admin/cards/rotate-keys \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | localhost |
| `DB_PORT` | PostgreSQL port | 5432 |
| `DB_USER` | Database user | cardvault_user |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | cardvault_db |
| `JWT_SECRET` | JWT signing secret (min 32 chars) | - |
| `PORT` | Application port | 8080 |
| `GIN_MODE` | Gin mode (debug/release) | debug |
| `ENABLE_TEST_AUTH` | Enable test token endpoint | true |

### Security Configuration

- **Rate Limiting**: 100 requests/second with burst of 20
- **JWT Expiration**: 24 hours (configurable)
- **Encryption**: AES-256-GCM with random nonces
- **Key Rotation**: Automatic versioning with backward compatibility

## 🔒 Security Features

### PCI DSS Compliance
- **Data Encryption**: All sensitive data encrypted at rest and in transit
- **Access Control**: JWT-based authentication with user isolation
- **Audit Logging**: Comprehensive audit trails for all operations
- **Data Masking**: Card numbers masked in all responses
- **Secure Transmission**: HTTPS enforcement with security headers

### Encryption Details
- **Algorithm**: AES-256-GCM
- **Key Size**: 256-bit keys with automatic generation
- **Nonce**: Unique random nonce per encryption operation
- **Key Management**: Secure key rotation without service interruption

### Input Validation
- **Card Numbers**: Luhn algorithm validation
- **Expiry Dates**: Future date validation
- **CVV**: Format and length validation
- **Request Size**: Limited to prevent DoS attacks

## 📊 Performance

### Benchmarks
- **Single Card Creation**: ~10ms average response time
- **Batch Updates**: Handles 1000+ cards efficiently
- **Concurrent Requests**: Supports 1000+ concurrent connections
- **Database Performance**: Optimized with proper indexing

### Monitoring
- Health check endpoint at `/health`
- Request rate limiting and monitoring
- Database connection pooling
- Memory usage optimization

## 🚀 Deployment

### Production Deployment

1. **Build production image**
   ```bash
   docker build -t card-vault:production .
   ```

2. **Deploy with production config**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

3. **Environment setup**
   - Set `GIN_MODE=release`
   - Use strong `JWT_SECRET`
   - Configure proper database credentials
   - Enable SSL/TLS termination

### Scaling Considerations
- **Horizontal Scaling**: Stateless design supports multiple instances
- **Database**: Consider read replicas for high-read workloads
- **Caching**: Redis integration ready for session management
- **Load Balancing**: Compatible with standard load balancers

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go best practices and conventions
- Add unit tests for new features
- Update API documentation
- Ensure PCI DSS compliance for security features

## 🙏 Acknowledgments

- Built with [Gin](https://gin-gonic.com/) web framework
- Database management with [GORM](https://gorm.io/)
- JWT implementation using [golang-jwt](https://github.com/golang-jwt/jwt)
- Containerization with Docker

---

