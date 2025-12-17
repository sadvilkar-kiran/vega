# Vega

**Speed Up Your Development**

Vega is a high-performance Go web framework designed for rapid application development. Built with speed and simplicity in mind.

## Features

- ⚡ Lightning-fast routing with Chi
- 🔐 Built-in authentication & session management
- 📧 Email support (SMTP + APIs)
- 💾 Multiple database support (PostgreSQL, MySQL)
- 🗄️ Caching (Redis, Badger)
- 📝 Template rendering (Jet, Go templates)
- 🛡️ Security features (CSRF, encryption)
- 🚀 CLI tool for scaffolding
- 📦 Database migrations
- ⏰ Cron job scheduler

## Installation

### Framework
```bash
go get github.com/sadvilkar-kiran/vega
```

### CLI Tool
```bash
go install github.com/sadvilkar-kiran/vega/cmd/vega@latest
```

After installation, you can use the `vega` command:
```bash
vega new myapp
vega help
```

## Quick Start

```go
package main

import (
    "log"
    "os"
    
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

## Documentation

Full documentation coming soon. For now, check the source code and examples.

## License

MIT License - see LICENSE file

## Author

**Kiran Sadvilkar**

- GitHub: [@sadvilkar-kiran](https://github.com/sadvilkar-kiran)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

