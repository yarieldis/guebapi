# Running guebapi with Docker + Kubernetes (API Gateway)

## Overview

This document analyzes how the guebapi project can be deployed using Docker and
Kubernetes instead of (or alongside) the current AWS API Gateway + Lambda setup.

The short answer: **yes, it is very doable**, because the application already
separates the HTTP layer (Gin) from the deployment layer.

## The key insight

`app.BuildRouter()` returns a plain `*gin.Engine`, which is a standard
`http.Handler`. The two entry points treat it differently:

| Entry point            | Backend  | How it is served                                          |
|------------------------|----------|-----------------------------------------------------------|
| `main.go` (local)      | memory   | `r.Run(":8080")` — plain HTTP server                      |
| `cmd/lambda/main.go`   | DynamoDB | `ginadapter.NewV2(r).ProxyWithContext` — converts API Gateway events → HTTP requests |

So the "API Gateway" in the current setup is **AWS API Gateway**, and the Lambda
adapter translates its events into HTTP for Gin. In Docker + Kubernetes, that
relationship is flipped: Gin runs as a **normal long-lived HTTP server inside a
container**, and the gateway role is filled by an **in-cluster ingress/gateway**.

## The mapping

```
AWS (today)                        →  Kubernetes (Docker)
─────────────────────────────────────────────────────────────────
API Gateway  (routes + $default)   →  Ingress / Gateway API
Lambda       (event → HTTP)        →  Deployment (pods running the Gin server)
ginadapter   (event translation)   →  (gone — no translation needed)
DynamoDB     (storage)             →  DynamoDB via IRSA, or a K8s-friendly DB
```

The gateway options that fill AWS API Gateway's role:

- **Kong**, **Traefik**, **NGINX Ingress**, **Emissary-Ingress (Ambassador)** — traditional ingress
- **Envoy Gateway**, **Istio**, **Linkerd** — service-mesh/gateway approach
- On EKS specifically: **AWS Load Balancer Controller** — an ALB in front of the
  cluster (closest analog to the current setup)

## What needs to change in the code

Only one small thing: a **unified server entry point** is needed that picks the
backend from `STORAGE_BACKEND`. Currently `main.go` hardcodes the memory
repository and `cmd/lambda/main.go` hardcodes DynamoDB — neither is right for a
container.

```go
// cmd/server/main.go
func main() {
    cfg := config.Load()

    var userRepo user.Repository
    switch cfg.Storage.Backend {
    case "dynamodb":
        awsCfg, _ := config.LoadDefaultConfig(context.Background(),
            config.WithRegion(cfg.AWS.Region))
        userRepo = user.NewDynamoDBRepository(dynamodb.NewFromConfig(awsCfg), cfg.Storage.DynamoDBTable)
    default:
        userRepo = user.NewMemoryRepository()
    }

    r := app.BuildRouter(cfg, userRepo)
    if err := r.Run(cfg.Server.Address()); err != nil {
        log.Fatalf("failed to start server: %v", err)
    }
}
```

`SERVER_HOST` defaults to `""`, so `Address()` returns `:8080`, which binds
`0.0.0.0` — correct for containers.

## Dockerfile (multi-stage)

```dockerfile
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/guebapi ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/guebapi /guebapi
EXPOSE 8080
ENTRYPOINT ["/guebapi"]
```

## Kubernetes (Deployment + Service + Ingress)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: guebapi
spec:
  replicas: 2
  selector:
    matchLabels: { app: guebapi }
  template:
    metadata:
      labels: { app: guebapi }
    spec:
      containers:
        - name: guebapi
          image: guebapi:latest
          ports: [{ containerPort: 8080 }]
          env:
            - name: STORAGE_BACKEND
              value: dynamodb          # or memory for a demo
            - name: JWT_SECRET_KEY
              valueFrom:
                secretKeyRef: { name: guebapi, key: jwt-secret }
          readinessProbe:
            httpGet: { path: /health, port: 8080 }   # /health is perfect for this
          livenessProbe:
            httpGet: { path: /health, port: 8080 }
---
apiVersion: v1
kind: Service
metadata:
  name: guebapi
spec:
  selector: { app: guebapi }
  ports: [{ port: 80, targetPort: 8080 }]
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: guebapi
spec:
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service: { name: guebapi, port: { number: 80 } }
```

## Important caveat: storage

The `memory` backend is a `sync.RWMutex` on a single process. In Kubernetes with
`replicas: 2`, each pod has **its own separate store**, and data is lost on pod
restart. For anything real, use `dynamodb` (with **IAM Roles for Service
Accounts / IRSA** on EKS) or swap in Postgres/RDS. The `user.Repository`
interface already makes that swap clean.

## Summary

- Nothing about the core application needs to change — it is already a plain HTTP server.
- The AWS API Gateway + Lambda adapter is replaced by an Ingress/Gateway + Deployment.
- One new entry point (`cmd/server`) is needed to respect `STORAGE_BACKEND`, plus
  a `Dockerfile` and Kubernetes manifests.
