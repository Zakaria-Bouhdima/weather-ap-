# WeatherApp

A microservices weather application deployed on Kubernetes with full observability stack.

## Architecture

```
User → Nginx (reverse proxy)
         ├── UI Service      (Node.js)
         ├── Auth Service    (Go / Gin)
         └── Weather Service (Python / Flask)

Monitoring: Prometheus + Grafana + Loki + Promtail
CI/CD:      GitHub Actions → GHCR
```

## Services

| Service | Language | Image |
|---------|----------|-------|
| nginx   | -        | `ghcr.io/<owner>/weatherapp-nginx` |
| ui      | Node.js  | `ghcr.io/<owner>/weatherapp-ui` |
| auth    | Go       | `ghcr.io/<owner>/weatherapp-auth` |
| weather | Python   | `ghcr.io/<owner>/weatherapp-weather` |

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Helm](https://helm.sh/docs/intro/install/)
- [Make](https://www.gnu.org/software/make/)

## Quick Start

```bash
# 1. Create local Kind cluster
make cluster-up

# 2. Configure secrets
make secrets          # follow printed instructions

# 3. Deploy application
make deploy

# 4. Deploy monitoring (Prometheus + Grafana + Loki)
make monitoring-install
```

Grafana is available at `http://localhost:30300` (default credentials: `admin` / `admin`).

## Makefile Targets

| Target               | Description                              |
|----------------------|------------------------------------------|
| `cluster-up`         | Create Kind cluster                      |
| `cluster-down`       | Delete Kind cluster                      |
| `secrets`            | Print secrets setup instructions         |
| `deploy`             | Install / upgrade all app Helm charts    |
| `destroy`            | Uninstall all app Helm charts            |
| `monitoring-install` | Deploy Prometheus, Grafana, Loki         |
| `monitoring-destroy` | Teardown monitoring stack                |

## Monitoring

The observability stack includes:

- **Prometheus** — metrics collection with custom alert rules
- **Grafana** — pre-built WeatherApp dashboard (request rate, error rate, p95 latency, pod health)
- **Loki + Promtail** — log aggregation from all pods
- **Alertmanager** — routes critical alerts to Slack

### SLOs

| SLO | Target |
|-----|--------|
| Error rate (5xx) | < 0.5% |
| p95 latency | < 500ms |

## CI/CD

GitHub Actions builds and pushes Docker images to GHCR on every push to `main`.

```
push to main
  ├── build-auth    → ghcr.io/<owner>/weatherapp-auth:latest
  ├── build-ui      → ghcr.io/<owner>/weatherapp-ui:latest
  ├── build-weather → ghcr.io/<owner>/weatherapp-weather:latest
  └── build-nginx   → ghcr.io/<owner>/weatherapp-nginx:latest
```

## Secrets

Copy the example files and fill in real values before deploying:

```bash
cp .env.example .env
cp k8s/secrets.yaml.example k8s/secrets.yaml
```

Never commit `.env` or `k8s/secrets.yaml` — they are in `.gitignore`.
