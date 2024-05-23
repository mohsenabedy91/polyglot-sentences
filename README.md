# Polyglot Sentences
Polyglot Sentences is a Go-based application designed to help users learn and master sentences in multiple languages. The app provides a wide range of sentence structures and vocabulary to facilitate language learning through practical and contextual examples.

# Go Version
- The project uses Go `1.22.3`

# Project Architecture
- The project architecture is based on the `Hexagonal Architecture`.
- The project is structured in a way that it is easy to understand and navigate through.


```tree-extended
📁polyglot-sentences/
├── 📁.github/
├── 📁cmd/
│   ├── 📁http/
│   │   └── 📄main.go
│   ├── 📁migration/
│   │   └── 📄main.go
│   └── 📁worker/
│       └── 📄main.go
├── 📁docs/
│   ├── 📄docs.go
│   ├── 📄swagger.json
│   └── 📄swagger.yaml
├── 📁internal/
│   ├── 📁adapter/
│   │   ├── 📁messagebroker/
│   │   │   └── 📁rabbitmq/
│   │   │       ├── 📄connection.go
│   │   │       ├── 📄producer.go
│   │   │       └── 📄consumer.go
│   │   ├── 📁http/
│   │   │   ├── 📁handler/
│   │   │   │   ├── 📄health.go
│   │   │   │   ├── 📄status_code_mapping.go
│   │   │   │   └── 📄user.go
│   │   │   ├── 📁middleware/
│   │   │   │   └── 📄custom_recovery.go
│   │   │   ├── 📁request/
│   │   │   │   └── 📄user.go
│   │   │   ├── 📁presenter/
│   │   │   │   ├── 📄base.go
│   │   │   │   └── 📄user.go
│   │   │   ├── 📁routes/
│   │   │   │   ├── 📄router.go
│   │   │   │   └── 📄swagger.go
│   │   │   └── 📁validations/
│   │   │       └── 📄validator.go
│   │   └── 📁storage/
│   │       ├── 📁postgres/
│   │       │   ├── 📁migrations/
│   │       │   │   ├── 📄202404031147_create_users_table.down.sql
│   │       │   │   └── 📄202404031147_create_users_table.up.sql
│   │       │   ├── 📁repository/
│   │       │   │   └── 📄user.go
│   │       │   └── 📄db.go
│   │       └── 📁redis/
│   │           └── 📄redis.go
│   └── 📁core/
│       ├── 📁config/
│       │   └── 📄config.go
│       ├── 📁domain/
│       │   ├── 📄base.go
│       │   └── 📄user.go
│       ├── 📁port/
│       │   ├── 📄message_broker.go
│       │   └── 📄user.go
│       └── 📁service/
│           └── 📁userservice/
│               └── 📄user.go
├── 📁logs/
│   └── 📄logs-2024-05-21.log
├── 📁pkg/
│   ├── 📁claim/
│   │   └── 📄gin.go
│   ├── 📁logger/
│   │   ├── 📄const.go
│   │   └── 📄logger.go
│   ├── 📁password/
│   │   └── 📄password.go
│   ├── 📁serviceerror/
│   │   ├── 📄error_message.go
│   │   └── 📄service_error.go
│   └── 📁translation/
│       ├── 📄trans.go
│       └── 📁lang/
│           ├── 📄ar.json
│           ├── 📄en.json
│           └── 📄fa.json
├── 📄go.mod
├── 📄.env
└── 📄docker-compose.yml
```

# Profiling
We use pprof tool for get CPU, go routine and memory leak

- [pprof](http://localhost:2526/debug/pprof/)
- [goroutine](http://localhost:2526/debug/pprof/goroutine?debug=1)

```bash
curl http://localhost:2526/debug/pprof/goroutine --output goroutine.o

go tool pprof -http=:2020 goroutine.o
```
if debug mode was true its work

# Requirements
#### Authentication/Authorization 
#### User Management
#### Questions Management
#### Questions Planner
#### Telegram Integration