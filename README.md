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
│   ├── 📁authserver/
│   │   └── 📄http.go
│   ├── 📁migration/
│   │   └── 📄main.go
│   └── 📁userserver/
│       ├── 📄grpc.go
│       └── 📄http.go
├── 📁deploy/
│   ├── 📄Deployment.yml
│   └── 📄Service.yml
├── 📁docker/
│   ├── 📁alertmanager/
│   ├── 📁elk/
│   ├── 📁grafana/
│   └── 📁prometheus/
├── 📁docs/
│   ├── 📄docs.go
│   ├── 📄swagger.json
│   └── 📄swagger.yaml
├── 📁internal/
│   ├── 📁adapter/
│   │   ├── 📁constant/
│   │   │   └── 📄messages.go
│   │   ├── 📁grpc/
│   │   │   ├── 📁client/
│   │   │   │   └── 📄user_client.go
│   │   │   ├── 📁proto/
│   │   │   │   ├── 📄user.pb.go
│   │   │   │   ├── 📄user.proto
│   │   │   │   └── 📄user_grpc.pb.go
│   │   │   └── 📁server/
│   │   │       └── 📄user_server.go
│   │   ├── 📁http/
│   │   │   ├── 📁handler/
│   │   │   │   ├── 📄health.go
│   │   │   │   ├── 📄status_code_mapping.go
│   │   │   │   └── 📄user.go
│   │   │   ├── 📁middleware/
│   │   │   │   └── 📄custom_recovery.go
│   │   │   ├── 📁presenter/
│   │   │   │   ├── 📄base.go
│   │   │   │   └── 📄user.go
│   │   │   ├── 📁request/
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
│   │           └── 📄db.go
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
├── 📁proto/
│   └── 📁common/
│       ├── 📄error_details.pb.go
│       └── 📄error_details.proto
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
## Authentication/Authorization:
- Proto buffer:
There we need to get user details for this matters you should run protoc command for user management service
```bash
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
proto/common/error_details.proto
```
```bash
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
internal/adapter/grpc/proto/user/user.proto
```
## User Management
## Questions Management
## Questions Planner
## Telegram Integration