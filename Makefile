CLUSTER_NAME    := weatherapp
NAMESPACE       := default
MONITORING_NS   := monitoring

.PHONY: cluster-up cluster-down deploy destroy secrets monitoring-install monitoring-destroy

cluster-up:
	kind create cluster --config kind/cluster.yaml
	kubectl wait --for=condition=Ready nodes --all --timeout=60s

cluster-down:
	kind delete cluster --name $(CLUSTER_NAME)

secrets:
	@echo "Apply secrets before deploying:"
	@echo "  cp k8s/secrets.yaml.example k8s/secrets.yaml"
	@echo "  # edit k8s/secrets.yaml with real values"
	@echo "  kubectl apply -f k8s/secrets.yaml"

deploy:
	helm upgrade --install weatherapp-auth  ./weatherapp-auth  -n $(NAMESPACE)
	helm upgrade --install weatherapp-weather ./weatherapp-weather -n $(NAMESPACE)
	helm upgrade --install weatherapp-ui     ./weatherapp-ui    -n $(NAMESPACE)
	helm upgrade --install weatherapp-nginx  ./weatherapp-nginx  -n $(NAMESPACE)

destroy:
	helm uninstall weatherapp-nginx   -n $(NAMESPACE) || true
	helm uninstall weatherapp-ui      -n $(NAMESPACE) || true
	helm uninstall weatherapp-weather -n $(NAMESPACE) || true
	helm uninstall weatherapp-auth    -n $(NAMESPACE) || true

monitoring-install:
	kubectl create namespace $(MONITORING_NS) --dry-run=client -o yaml | kubectl apply -f -
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo update
	kubectl create configmap weatherapp-dashboard \
		--from-file=weatherapp-dashboard.json=monitoring/dashboards/weatherapp-dashboard.json \
		-n $(MONITORING_NS) --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f monitoring/rules/weatherapp-alerts.yaml
	helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
		-n $(MONITORING_NS) -f monitoring/kube-prometheus-stack-values.yaml
	helm upgrade --install loki grafana/loki \
		-n $(MONITORING_NS) -f monitoring/loki-values.yaml

monitoring-destroy:
	helm uninstall loki      -n $(MONITORING_NS) || true
	helm uninstall prometheus -n $(MONITORING_NS) || true
	kubectl delete namespace $(MONITORING_NS) || true
