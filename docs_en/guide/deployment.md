# Deployment

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Docker

#### Building the Image

```bash
make docker-build
```

Or manually:

```bash
docker build -t mon-projet:latest .
```

The generated Dockerfile uses a multi-stage build:

```dockerfile
# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Stage 2: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY .env.example .env

EXPOSE 8080
CMD ["./main"]
```

**Advantages**:
- Lightweight final image (~15-20MB vs ~1GB)
- Security (minimal alpine image)
- Static binary (no dependencies)

#### Run with Docker

```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PASSWORD=postgres \
  -e JWT_SECRET=<votre_secret> \
  mon-projet:latest
```

**Note**: `host.docker.internal` allows access to localhost from Docker.

#### Docker Compose

If `docker-compose.yml` is generated:

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: mon-projet
      JWT_SECRET: ${JWT_SECRET}
    depends_on:
      - postgres

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: mon-projet
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

**Launch**:

```bash
# Set JWT_SECRET
export JWT_SECRET=$(openssl rand -base64 32)

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop
docker-compose down
```

### Kubernetes

Basic manifests for K8s deployment:

#### Secret

```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: mon-projet-secret
type: Opaque
stringData:
  jwt-secret: "<votre_secret_base64>"
  db-password: "postgres"
```

```bash
kubectl apply -f k8s/secret.yaml
```

#### Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mon-projet
  labels:
    app: mon-projet
spec:
  replicas: 3
  selector:
    matchLabels:
      app: mon-projet
  template:
    metadata:
      labels:
        app: mon-projet
    spec:
      containers:
      - name: mon-projet
        image: mon-projet:latest
        ports:
        - containerPort: 8080
        env:
        - name: APP_PORT
          value: "8080"
        - name: DB_HOST
          value: postgres-service
        - name: DB_USER
          value: postgres
        - name: DB_NAME
          value: mon-projet
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mon-projet-secret
              key: db-password
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: mon-projet-secret
              key: jwt-secret
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
```

```bash
kubectl apply -f k8s/deployment.yaml
```

#### Service

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: mon-projet-service
spec:
  selector:
    app: mon-projet
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
  type: LoadBalancer
```

```bash
kubectl apply -f k8s/service.yaml
```

#### Deploy PostgreSQL (StatefulSet)

```yaml
# k8s/postgres.yaml
apiVersion: v1
kind: Service
metadata:
  name: postgres-service
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
  clusterIP: None
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
spec:
  serviceName: postgres-service
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:16-alpine
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_DB
          value: mon-projet
        - name: POSTGRES_USER
          value: postgres
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mon-projet-secret
              key: db-password
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
  - metadata:
      name: postgres-storage
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 10Gi
```

```bash
kubectl apply -f k8s/postgres.yaml
```

### CI/CD with GitHub Actions

The generated workflow (`.github/workflows/ci.yml`):

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.25'

    - name: golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest

  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_DB: test_db
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.25'

    - name: Run tests
      env:
        DB_HOST: localhost
        DB_PORT: 5432
        DB_USER: postgres
        DB_PASSWORD: postgres
        DB_NAME: test_db
        JWT_SECRET: test-secret
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out

  build:
    runs-on: ubuntu-latest
    needs: [quality, test]
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.25'

    - name: Build
      run: go build -v -o mon-projet cmd/main.go
```

**Pipeline**:
1. **Quality**: golangci-lint
2. **Test**: Tests with PostgreSQL (service container)
3. **Build**: Build verification

### Production Deployment

#### Pre-deployment Checklist

- [ ] All tests pass (`make test`)
- [ ] Lint passes (`make lint`)
- [ ] Environment variables configured
- [ ] JWT_SECRET generated (strong, random)
- [ ] DB_SSLMODE=require
- [ ] DB migrations executed
- [ ] Health check working
- [ ] Logs configured
- [ ] Monitoring in place

#### Recommended Platforms

**1. Google Cloud Run** (simplest):

```bash
# Build and push image
gcloud builds submit --tag gcr.io/PROJECT_ID/mon-projet

# Deploy
gcloud run deploy mon-projet \
  --image gcr.io/PROJECT_ID/mon-projet \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars JWT_SECRET=$JWT_SECRET,DB_HOST=$DB_HOST
```

**2. AWS ECS/Fargate**:

- Build image → Push to ECR
- Create Task Definition
- Create ECS Service
- Configure ALB

**3. Heroku**:

```bash
# Login
heroku login

# Create app
heroku create mon-projet

# Add PostgreSQL
heroku addons:create heroku-postgresql:hobby-dev

# Set env vars
heroku config:set JWT_SECRET=$(openssl rand -base64 32)

# Deploy
git push heroku main
```

**4. Kubernetes** (most flexible):

```bash
# Apply all manifests
kubectl apply -f k8s/

# Check status
kubectl get pods
kubectl get services

# View logs
kubectl logs -f deployment/mon-projet
```

---


---

## Navigation

**Previous**: [Security](security.md)  
**Next**: [Monitoring](monitoring.md)  
**Index**: [Guide Index](index.md)
