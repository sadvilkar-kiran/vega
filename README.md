# Vega Framework

**Speed Up Your Development** 🚀

Vega is a high-performance Go web framework designed for rapid application development. Built with speed and simplicity in mind, Vega helps you build production-ready web applications faster.

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/sadvilkar-kiran/vega.svg?style=flat-square)](https://github.com/sadvilkar-kiran/vega/releases)

## Features

- ⚡ **Lightning-fast routing** with Chi router
- 🔐 **Built-in authentication** & session management
- 📧 **Email support** (SMTP + Mailgun, SendGrid, SparkPost APIs)
- 💾 **Multiple database support** (PostgreSQL, MySQL/MariaDB)
- 🗄️ **Caching** (Redis, Badger)
- 📝 **Template rendering** (Jet templates, Go standard templates)
- 🛡️ **Security features** (CSRF protection, AES encryption, URL signing)
- 🚀 **CLI tool** for rapid scaffolding
- 📦 **Database migrations** built-in
- ⏰ **Cron job scheduler**
- ✅ **Form validation** with govalidator
- 📊 **JSON/XML response helpers**

## Installation

### Step 1: Install the CLI Tool

The easiest way to get started with Vega is using the CLI tool:

```bash
go install github.com/sadvilkar-kiran/vega/cmd/vega@latest
```

**Note for Windows users:** After installation, the CLI will be at `C:\Users\YourName\go\bin\vega.exe`. You can either:
- Use the full path: `C:\Users\YourName\go\bin\vega.exe`
- Add `C:\Users\YourName\go\bin` to your PATH environment variable

**Verify installation:**
```bash
vega version
# Should output: Vega CLI version: 1.0.0
```

### Step 2: Create Your First App

```bash
# Create a new Vega application
vega new myapp

# Navigate into your app
cd myapp

# Install dependencies
go mod tidy

# Run your app
go run .
```

Your app will be available at **http://localhost:4000**

## Quick Start

### Using the CLI (Recommended)

```bash
# 1. Install CLI
go install github.com/sadvilkar-kiran/vega/cmd/vega@latest

# 2. Create new app
vega new myapp

# 3. Enter app directory
cd myapp

# 4. Install dependencies
go mod tidy

# 5. Run the app
go run .
```

That's it! Your app is running with:
- ✅ Complete project structure
- ✅ All necessary files generated
- ✅ View templates included
- ✅ `.env` file with auto-generated encryption key
- ✅ Example routes and handlers

### Manual Installation

If you prefer to set up manually:

```bash
# Install the framework
go get github.com/sadvilkar-kiran/vega
```

Then in your code:

```go
package main

import (
    "log"
    "os"
    "net/http"
    
    "github.com/sadvilkar-kiran/vega"
)

func main() {
    path, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    // Initialize Vega
    v := &vega.Vega{}
    err = v.New(path)
    if err != nil {
        log.Fatal(err)
    }

    v.AppName = "myapp"

    // Add your routes
    v.Routes.Get("/", func(w http.ResponseWriter, r *http.Request) {
        v.Render.Page(w, r, "home", nil, nil)
    })

    // Start server
    v.ListenAndServe()
}
```

## CLI Commands

Vega CLI provides powerful commands to speed up development:

### Create New Application
```bash
vega new <app-name>
```
Creates a complete application structure with all necessary files and templates.

### Generate Code

```bash
# Generate encryption key
vega make key

# Create database migration
vega make migration <migration-name>

# Create handler
vega make handler <handler-name>

# Create model
vega make model <model-name>

# Scaffold authentication (users, tokens, middleware)
vega make auth

# Create session table
vega make session

# Create mail templates
vega make mail <template-name>
```

### Database Migrations

```bash
# Run all pending migrations
vega migrate

# Rollback last migration
vega migrate down

# Rollback all migrations
vega migrate down all

# Reset database (down all, then up all)
vega migrate reset
```

### Help

```bash
vega help
vega version
```

## Project Structure

When you create a new app with `vega new myapp`, you get:

```
myapp/
├── main.go              # Application entry point
├── init-vega.go         # Framework initialization
├── routes.go            # Route definitions
├── convenience.go       # Helper functions
├── .env                 # Environment configuration
├── go.mod               # Go module file
├── handlers/            # HTTP handlers
│   └── handlers.go
├── middleware/          # Custom middleware
│   └── middleware.go
├── data/                # Data models
│   └── models.go
├── views/               # Template files
│   ├── home.jet
│   ├── layouts/
│   └── ...
├── migrations/          # Database migrations
├── mail/                # Email templates
├── public/              # Static files
└── logs/                # Application logs
```

## Configuration

Edit the `.env` file to configure your application:

```env
# Application
APP_NAME=myapp
APP_URL=http://localhost:4000
DEBUG=true
PORT=4000

# Database (PostgreSQL or MySQL)
DATABASE_TYPE=postgres
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=your_user
DATABASE_PASS=your_password
DATABASE_NAME=your_database

# Session Store (cookie, redis, mysql, postgres)
SESSION_TYPE=cookie

# Cache (redis, badger)
CACHE=

# Template Engine (go, jet)
RENDERER=jet

# Encryption Key (auto-generated, 32 characters)
KEY=your-32-character-encryption-key
```

## Examples

### JSON Response

```go
func (h *Handlers) MyHandler(w http.ResponseWriter, r *http.Request) {
    payload := map[string]interface{}{
        "message": "Hello from Vega!",
        "status": "success",
    }
    h.App.WriteJSON(w, http.StatusOK, payload)
}
```

### Template Rendering

```go
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
    td := &render.TemplateData{
        StringMap: map[string]string{
            "name": "Vega",
        },
    }
    h.App.Render.Page(w, r, "home", nil, td)
}
```

### Database Operations

```go
// Using upper/db
users, err := h.Models.Users.GetAll()
user, err := h.Models.Users.Get(1)
id, err := h.Models.Users.Insert(newUser)
```

### Caching

```go
// Set cache
h.App.Cache.Set("key", "value")

// Get cache
value, err := h.App.Cache.Get("key")

// Delete cache
h.App.Cache.Forget("key")
```

## Documentation

- **Full Documentation**: Coming soon
- **Examples**: Check the [vegaapp](https://github.com/sadvilkar-kiran/vegaapp) repository
- **API Reference**: [GoDoc](https://pkg.go.dev/github.com/sadvilkar-kiran/vega)

## Requirements

- Go 1.21 or higher
- (Optional) PostgreSQL or MySQL for database features
- (Optional) Redis for caching/sessions

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Author

**Kiran Sadvilkar**

- GitHub: [@sadvilkar-kiran](https://github.com/sadvilkar-kiran)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Support

- **Issues**: [GitHub Issues](https://github.com/sadvilkar-kiran/vega/issues)
- **Discussions**: [GitHub Discussions](https://github.com/sadvilkar-kiran/vega/discussions)

## Roadmap

- [ ] Additional documentation and examples
- [ ] More database drivers
- [ ] WebSocket support
- [ ] GraphQL support
- [ ] Additional middleware

---

**Built with ❤️ using Go**

**Vega - Speed Up Your Development** 🚀
