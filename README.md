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
│   ├── 📁apigateway/
│   │   └── 📄main.go
│   ├── 📁authserver/
│   │   └── 📄main.go
│   ├── 📁migration/
│   │   └── 📄main.go
│   ├── 📁setup/
│   │   └── 📄setup.go
│   ├── 📁userserver/
│   │   └── 📄main.go
│   └── 📁worker/
│       └── 📄main.go
├── 📁deploy/
│   ├── 📄Deployment.yml
│   └── 📄Service.yml
├── 📁docker/
│   ├── 📁alertmanager/
│   ├── 📁elk/
│   ├── 📁grafana/
│   ├── 📁kong/
│   └── 📁prometheus/
├── 📁docs/
│   ├── 📄docs.go
│   ├── 📄swagger.json
│   └── 📄swagger.yaml
├── 📁internal/
│   ├── 📁adapter/
│   │   ├── 📁constant/
│   │   │   └── 📄messages.go
│   │   ├── 📁email/
│   │   │   ├── 📁mocks/
│   │   │   │   └── 📄mock_sendgrid.go
│   │   │   ├── 📄sendgrid.go
│   │   │   └── 📄sendgrid_test.go
│   │   ├── 📁grpc/
│   │   │   ├── 📁client/
│   │   │   │   └── 📄user_client.go
│   │   │   ├── 📁proto/
│   │   │   │   └── 📁user/
│   │   │   │       ├── 📄user.pb.go
│   │   │   │       ├── 📄user.proto
│   │   │   │       └── 📄user_grpc.pb.go
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
│   │   │   │   ├── 📄base.go
│   │   │   │   └── 📄user.go
│   │   │   ├── 📁routes/
│   │   │   │   ├── 📄router.go
│   │   │   │   └── 📄swagger.go
│   │   │   └── 📁validations/
│   │   │       └── 📄validator.go
│   │   ├── 📁messagebroker/
│   │   │   ├── 📄queue.go
│   │   │   └── 📄rabbitmq.go
│   │   ├── 📁minio/
│   │   │   └── 📄client.go
│   │   └── 📁storage/
│   │       ├── 📁postgres/
│   │       │   ├── 📁migrations/
│   │       │   │   ├── 📄202404031147_create_users_table.down.sql
│   │       │   │   └── 📄202404031147_create_users_table.up.sql
│   │       │   ├── 📁authrepository/
│   │       │   │   ├── 📄access_control.go
│   │       │   │   ├── 📄postgres_test.go
│   │       │   │   ├── 📄permission.go
│   │       │   │   ├── 📄role.go
│   │       │   │   └── 📄unit_of_work.go
│   │       │   ├── 📁userrepository/
│   │       │   │   ├── 📄postgres_test.go
│   │       │   │   ├── 📄unit_of_work.go
│   │       │   │   ├── 📄user.go
│   │       │   │   └── 📄user_test.go
│   │       │   └── 📄db.go
│   │       └── 📁redis/
│   │           ├── 📁authrepository/
│   │           │   ├── 📄auth.go
│   │           │   ├── 📄otp.go
│   │           │   └── 📄role.go
│   │           └── 📄db.go
│   └── 📁core/
│       ├── 📁config/
│       │   └── 📄config.go
│       ├── 📁constant/
│       │   └── 📄cache.go
│       ├── 📁domain/
│       │   ├── 📄access_control.go
│       │   ├── 📄base.go
│       │   ├── 📄grammer.go
│       │   ├── 📄language.go
│       │   ├── 📄permission.go
│       │   ├── 📄role.go
│       │   ├── 📄sentence.go
│       │   └── 📄user.go
│       ├── 📁event/
│       │   └── 📁authevent/
│       │       ├── 📄send_email_otp_queue.go
│       │       ├── 📄send_reset_password_link_queue.go
│       │       └── 📄send_welcome_queue.go
│       ├── 📁port/
│       │   ├── 📄access_control.go
│       │   ├── 📄aut.go
│       │   ├── 📄email.go
│       │   ├── 📄event.go
│       │   ├── 📄otp.go
│       │   ├── 📄permission.go
│       │   ├── 📄role.go
│       │   └── 📄user.go
│       ├── 📁service/
│       │   ├── 📁authservice/
│       │   │   └── 📄jwt.go
│       │   ├── 📁roleservice/
│       │   │   ├── 📄cache.go
│       │   │   └── 📄role.go
│       │   └── 📁userservice/
│       │       └── 📄user.go
│       └── 📁views/
│           └── 📁email/
│               ├── 📁auth
│               │   ├── 📄verify_email.html
│               │   └── 📄welcome.html
│               └── 📄base.html
├── 📁logs/
│   └── 📄logs-2024-05-21.log
├── 📁pkg/
│   ├── 📁claim/
│   │   └── 📄gin.go
│   ├── 📁helper/
│   │   ├── 📄authenticate.go
│   │   ├── 📄authenticate_bench_test.go
│   │   └── 📄string.go
│   ├── 📁logger/
│   │   ├── 📄const.go
│   │   └── 📄logger.go
│   ├── 📁metrics/
│   │   ├── 📄counters.go
│   │   └── 📄histograms.go
│   ├── 📁oauth/
│   │   └── 📄google.go
│   ├── 📁serviceerror/
│   │   ├── 📄error_message.go
│   │   ├── 📄grpc.go
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

# API Gateway
There is an API gateway, and we used of `Kong` for management.
The APIs available on `http://localhost:8000` and for Dashboard you can open with this link:
[workspaces](http://localhost:8002/default/overview)

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
protoc --experimental_allow_proto3_optional --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
internal/adapter/grpc/proto/user/user.proto
```
## User Management
## Questions Management
## Questions Planner
## Telegram Integration