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

# Cloud

## Install minikube in local machine
1. Download the Minikube binary:
```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
```
2. Install Minikube and remove the binary:
```bash
sudo install minikube-linux-amd64 /usr/local/bin/minikube
```
```bash
rm minikube-linux-amd64
```
3. Start Minikube:
```bash
minikube start
```
```bash
minikube start --memory=6144 --cpus=4
```
4. Verify the Minikube installation:
```bash
minikube status
```
5. Check Kubernetes pods across all namespaces:
```bash
kubectl get po -A
```
6. Configure Docker to use Minikube's Docker daemon:
```bash
eval $(minikube -p minikube docker-env)
```
7. List Docker images:
```bash
docker images
```

### Create Polyglot sentences namespace
1. Create the `polyglot-sentences` namespace:
```bash
kubectl create namespace polyglot-sentences
```
2. Set the current context to the `polyglot-sentences` namespace:
```bash
kubectl config set-context --current --namespace=polyglot-sentences
```

### Create config map
1. Apply the configuration map from the specified YAML file:
```bash
kubectl apply -f deploy/configs/config-maps.yaml
```
2. Verify the config maps:
```bash
kubectl get configmaps
```
```bash
kubectl describe configmap polyglot-sentences-env-config
```
```bash
kubectl describe configmap polyglot-sentences-file-config
```

### Create Services
1. Apply the `auth-service` configuration:
```bash
kubectl apply -f deploy/authservice/service.yaml
```
2. Apply the gRPC user management service configuration:
```bash
kubectl apply -f deploy/userservice/grpc-service.yaml
```
3. Apply the HTTP user management service configuration:
```bash
kubectl apply -f deploy/userservice/http-service.yaml
```
4. Verify the created services:
```bash
kubectl get services -o wide
```

### Create Secrets for environments
1. Create a secret for the `polyglot-sentences` namespace:
```bash
kubectl -n=polyglot-sentences create secret generic polyglot-sentences-secret --from-literal JWT_ACCESS_TOKEN_SECRET="your-access-token-secret" --from-literal SEND_GRID_KEY="send-grid-key"
```
2. Verify the secret:
```bash
kubectl get secret polyglot-sentences-secret -o yaml
```

### Create deployments
1. Apply the user management deployment:
```bash
kubectl apply -f deploy/userservice/deployment.yaml
```
2. Apply the authentication deployment:
```bash
kubectl apply -f deploy/authservice/deployment.yaml
```
3. Apply the notification deployment:
```bash
kubectl apply -f deploy/notificationservice/deployment.yaml
```
4. Verify the deployments:
```bash
kubectl get deployments -o wide
```
5. Check the status of the pods:
```bash
kubectl get pods -o wide
```

> Apply all Micro services related configurations at once:
```bash
kubectl apply -f deploy/userservice
kubectl apply -f deploy/authservice
kubectl apply -f deploy/notificationservice
```

## Rollout deployments for apply new version images
1. Rollout restart for all deployments in the polyglot-sentences namespace:
```bash
kubectl rollout restart deployment -n polyglot-sentences
```
2. Rollout restart specific deployments:
```bash
kubectl rollout restart deployment.apps/auth-deployment -n polyglot-sentences
```
```bash
kubectl rollout restart deployment.apps/user-management-deployment -n polyglot-sentences
```
```bash
kubectl rollout restart deployment.apps/notification-deployment -n polyglot-sentences
```

# Jenkins
1. Create the `jenkins` namespace:
```bash
kubectl create namespace jenkins
```
2. Set the current context to the `jenkins` namespace:
```bash
kubectl config set-context --current --namespace=jenkins
```
3. Apply the Jenkins persistent volume configuration:
```bash
kubectl apply -f deploy/jenkins/persistent-volume.yaml
```
4. Apply the Jenkins persistent volume claim configuration:
```bash
kubectl apply -f deploy/jenkins/persistent-volume-claim.yaml
```
5. Verify the persistent volume:
```bash
kubectl get pvc
```
```bash
kubectl describe pvc jenkins-volume-claim
```
6. Apply the Jenkins service configuration:
```bash
kubectl apply -f deploy/jenkins/service.yaml
```
7. Apply the Jenkins master deployment configuration:
```bash
kubectl apply -f deploy/jenkins/master-deployment.yaml
```
8. Apply all Jenkins-related configurations at once:
```bash
kubectl apply -f deploy/jenkins
```

### Setup docker sock API
1. Edit the Docker service file:
```bash
sudo nano /lib/systemd/system/docker.service
```
2. Find and remove the following line:
`ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock`
3. Replace it with:
```
ExecStart=/usr/bin/dockerd -H tcp://0.0.0.0:4243 -H unix:///var/run/docker.sock
```
4. Reload the systemd daemon and restart Docker:
```bash
sudo systemctl daemon-reload
sudo service docker restart
```
5. Check the Docker service status:
```bash
sudo service docker status
```
6. Verify the Docker daemon is working:
```bash
curl http://localhost:4243/version
```

> Dashboard URL
```
http://jenkins.local:30080
```
- Get jenkins secrets
```bash
kubectl -n jenkins exec -it $(kubectl get pods -n jenkins -o jsonpath="{.items[0].metadata.name}") -- cat /var/jenkins_home/secrets/initialAdminPassword
```

### plugins
After running the Jenkins service, navigate to: `Manage Jenkins` -> `System Configuration` -> `Plugins` -> `Available plugins`. Search for and install the following plugins:
- **Kubernetes**
- **SSH Agent**
- **Blue Ocean**
- **ThinBackup**
- **Slack Notification**
- **Role-based Authorization Strategy**

## Setup Kubernetes plugins
To configure the Kubernetes plugin, follow these steps:
1. Go to: `Manage Jenkins` -> `Clouds` -> `New cloud`, select `Kubernetes`, and enter a name (e.g., `kubernetes`).
2. Set the `Kubernetes URL` by retrieving it with the following command:
```bash
kubectl cluster-info
```
3. Check the `Disable https certificate check` option.
4. Set the `Kubernetes Namespace` to jenkins.
5. Add credentials:
  - Select `Secret text`
    - ID: `JENKINS_SECRET`
    - For the secret, retrieve the Kubernetes service account token by following these steps:
```bash
kubectl create serviceaccount jenkins --namespace=jenkins
kubectl apply -f deploy/jenkins/token.yaml
kubectl create rolebinding jenkins-admin-binding --clusterrole=admin --serviceaccount=jenkins:jenkins --namespace=jenkins
TOKEN_NAME=$(kubectl get secret --namespace=jenkins | grep jenkins-token | awk '{print $1}')
kubectl describe secret $TOKEN_NAME --namespace=jenkins
```
6. Check `WebSocket` under the connection options.

## Credentials
- **Docker Hub**: Use `Username and Password`.
  - **ID**: `docker-hub-credentials`
  - **Username**: Your Docker Hub Username.
  - **Password**: Your Docker Hub Password.
- **Kubernetes**: Use service account token as `Secret text`.
- **GitHub App**: Use the following:
  - **ID**: `GitHub-APP`
  - **App ID**: `Your GitHub App ID`
  - **Token**: Convert and provide the token with `$ cat path/to/converted-github-app.pem`
- **DB_PASSWORD**: Use `Secret text` (Your test DB password).
- **SSH Agent**: Use `SSH Username with private key`:
  - **ID**: `k8s`
  - **Username**: Kubernetes host user
  - **Private key**:
    - **Generate**: `$ ssh-keygen -t rsa -b 4096 -C "jenkins@example.com"`
    - **Copy**: `$ ssh-copy-id «kubernetes host user»@«kubernetes remote address»`
    - **Retrieve value**: `$ cat ~/.ssh/id_rsa`


## Variable
- **DB_HOST**: Your DB Host address
- **DB_PORT**: 5425
- **DB_NAME**: Your test DB name
- **DB_USERNAME**: Your test DB username
- **REDIS_HOST**: Your Redis Host address
- **REDIS_PORT**: 6325
- **K8S_USER**: Kubernetes host user
- **K8S_REMOTE_ADDRESS**: Kubernetes remote address

## Jobs
### First Job: Polyglot Sentences Linting and Run Test
1. **Name**: `Polyglot Sentences linting and run test`
2. **Trigger**: Check `GitHub hook trigger for GITScm polling`
3. **Pipeline**:
   - **Definition**: `Pipeline script from SCM` 
   - **SCM**: `Git` 
   - **Repository URL**: `https://github.com/mohsenabedy91/polyglot-sentences.git`
   - **Credentials**: `GitHub-APP`
   - **Branches to build**:
     - **Branch Specifier**: `:^(?!origin/master$|origin/develop$).*`
   - **Script Path**: `jenkinsfile-linter-and-test`

### Second Job: Polyglot Sentences Deploy to Develop
1. **Name**: `Polyglot Sentences deploy to develop`
2. **Trigger**: `Check GitHub hook trigger for GITScm polling`
3. **Pipeline**:
   - **Definition**: `Pipeline script from SCM`
   - **SCM**: `Git`
   - **Repository URL**: `https://github.com/mohsenabedy91/polyglot-sentences.git`
   - **Credentials**: `GitHub-APP`
   - **Branches to build**:
     - **Branch Specifier**: `*/develop`
   - **Script Path**: `Jenkinsfile`

# Kong Api gateway
1. Create the `kong` namespace:
```bash
kubectl create namespace kong
```
2. Set the current context to the `kong` namespace:
```bash
kubectl config set-context --current --namespace=kong
```
3. Create a ConfigMap for `Kong` plugins:
```bash
kubectl create configmap kong-plugins --from-file=/path/to/polyglot-sentences/docker/kong/plugins/ps-authorize/
```
4. Verify the ConfigMap:
```bash
kubectl get configmaps
```
```bash
kubectl describe configmap kong-plugins
```
5. Create a secret for the Kong database:
```bash
kubectl -n=kong create secret generic kong-db-secrets --from-literal POSTGRES_PASSWORD="password"
```
6. Verify the secret:
```bash
kubectl get secret
```
```bash
kubectl describe secret kong-db-secrets
```
7. Apply the Kong persistent volume claim configuration:
```bash
kubectl apply -f deploy/kong/persistent-volume-claim.yaml
```
8. Verify the persistent volume:
```bash
kubectl get pvc
```
```bash
kubectl describe pvc gateway-postgres-volume-claim
```
9. Apply the PostgreSQL service configuration for Kong:
```bash
kubectl apply -f deploy/kong/postgres-service.yaml
```
10. Apply the PostgreSQL deployment configuration for Kong:
```bash
kubectl apply -f deploy/kong/postgres-deployment.yaml
```
11. Apply the Kong service configuration:
```bash
kubectl apply -f deploy/kong/kong-service.yaml
```
12. Apply the Kong ingress configuration:
```bash
kubectl apply -f deploy/kong/ingress.yaml
```
13. Apply the Kong deployment configuration:
```bash
kubectl apply -f deploy/kong/kong-deployment.yaml
```
14. Check the status of the pods:
```bash
kubectl get pods -o wide
```
15. Apply all Kong-related configurations at once:
```bash
kubectl apply -f deploy/kong
```
16. Rollout restart for all deployments in the `kong` namespace:
```bash
kubectl rollout restart deployment -n kong
```
> Dashboard URL
```
http://kong.local:30080
```
