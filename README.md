# Polyglot Sentences
Polyglot Sentences is a Go-based application designed to help users learn and master sentences in multiple languages. The app provides a wide range of sentence structures and vocabulary to facilitate language learning through practical and contextual examples.

# Go Version
- The project uses Go version `1.22.5`

# Installation and Setup

To simplify the installation and setup process for developers, we have provided an `install_service.sh` script. This script will:

- Check if Docker is installed, and if not, it will install Docker.
- Install Docker Compose if it is not already installed.
- Set up the environment file.
- Run Docker Compose to start the necessary containers.
- Wait for the PostgreSQL and Kong containers to be ready.
- Run the database migrations.
- Set up the services in the API Gateway.
- Set up the API Gateway.

### Steps to Install and Run the Service

- Make the script executable:

```bash 
chmod +x install_service.sh
```

- Run the script:

```bash
./install_service.sh
```

# Default admin user detail

- username(email): `john.doe@gmail.com`
- password: `QWer123!@#`

# Project Architecture
The project is structured using the Hexagonal Architecture. Here’s an overview of the directory structure:

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
│   ├── 📁notificationserver/
│   │   └── 📄main.go
│   ├── 📁setup/
│   │   └── 📄setup.go
│   └── 📁userserver/
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
│   │   │   ├── 📄mock_sendgrid.go
│   │   │   ├── 📄sendgrid.go
│   │   │   └── 📄sendgrid_test.go
│   │   ├── 📁grpc/
│   │   │   ├── 📁client/
│   │   │   │   ├── 📄mock_user_client.go
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
│   │   │   │   ├── 📄base_test.go
│   │   │   │   ├── 📄user.go
│   │   │   │   └── 📄user_test.go
│   │   │   ├── 📁request/
│   │   │   │   ├── 📄base.go
│   │   │   │   ├── 📄user.go
│   │   │   │   └── 📄user_test.go
│   │   │   ├── 📁routes/
│   │   │   │   ├── 📄auth_router.go
│   │   │   │   ├── 📄router.go
│   │   │   │   ├── 📄swagger.go
│   │   │   │   └── 📄user_router.go
│   │   │   └── 📁validations/
│   │   │       ├── 📄validator.go
│   │   │       └── 📄validator_test.go
│   │   ├── 📁messagebroker/
│   │   │   ├── 📄queue.go
│   │   │   └── 📄rabbitmq.go
│   │   ├── 📁minio/
│   │   │   └── 📄client.go
│   │   └── 📁storage/
│   │       ├── 📁postgres/
│   │       │   ├── 📁authrepository/
│   │       │   │   ├── 📄access_control.go
│   │       │   │   ├── 📄mock_access_control.go
│   │       │   │   ├── 📄mock_permission.go
│   │       │   │   ├── 📄mock_role.go
│   │       │   │   ├── 📄mock_unit_of_work.go
│   │       │   │   ├── 📄permission.go
│   │       │   │   ├── 📄role.go
│   │       │   │   └── 📄unit_of_work.go
│   │       │   ├── 📁migrations/
│   │       │   │   ├── 📄202404031147_create_users_table.down.sql
│   │       │   │   └── 📄202404031147_create_users_table.up.sql
│   │       │   ├── 📁tests/
│   │       │   │   ├── 📄access_control_test.go
│   │       │   │   ├── 📄permission_test.go
│   │       │   │   ├── 📄repositories_test.go
│   │       │   │   ├── 📄role_test.go
│   │       │   │   └── 📄user_test.go
│   │       │   ├── 📁userrepository/
│   │       │   │   ├── 📄mock_unit_of_work.go
│   │       │   │   ├── 📄mock_user.go
│   │       │   │   ├── 📄unit_of_work.go
│   │       │   │   └── 📄user.go
│   │       │   └── 📄db.go
│   │       └── 📁redis/
│   │           ├── 📁authrepository/
│   │           │   ├── 📄auth.go
│   │           │   ├── 📄mock_auth.go
│   │           │   ├── 📄mock_otp.go
│   │           │   ├── 📄mock_role.go
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
To profile the application, we use the `pprof` tool for CPU, goroutine, and memory usage data.

- [pprof](http://localhost:2526/debug/pprof/)
- [goroutine](http://localhost:2526/debug/pprof/goroutine?debug=1)

```bash
curl http://localhost:2526/debug/pprof/goroutine --output goroutine.o
go tool pprof -http=:2020 goroutine.o
```
Make sure the debug mode is enabled for the above links to work.

# API Gateway
We use `Kong` as the API gateway for managing the APIs. The APIs are available at `http://localhost:8000`. You can access the Kong Dashboard here:
[workspaces](http://localhost:8002/default/overview)

# Requirements
## Authentication/Authorization:
- Proto buffer:
  To get user details, run the protoc command for the user management service:
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

# Test
### Setup
Create a .env.test file with the necessary test environment variables.

### Running Tests
To run the tests, execute:
```bash
go test -cover -count=1 ./...
```
To generate a coverage profile:
```bash
go test -covermode=count -coverprofile=prof.out ./...
```
To visualize the coverage profile:
```bash
go tool cover -html=prof.out
```

### Test Suites
The project includes comprehensive test suites to ensure the functionality and reliability of the codebase, covering various components and features.

### Running All Tests
To run all the tests in the project:
```bash
go test ./... -v
```

# Linter
To install the linter package:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

To check code with the linter:
```bash
golangci-lint run
```
