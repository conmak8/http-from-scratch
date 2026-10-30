IMAGE_NAME := http-server
TAG := local
CLUSTER_NAME := observability

.PHONY: image
image:
	@echo "🐳 Building Docker image..."
	docker build -t $(IMAGE_NAME):$(TAG) .

.PHONY: deploy-local
deploy-local: image
	@echo "📦 Loading image into KIND cluster..."
	kind load docker-image $(IMAGE_NAME):$(TAG) --name $(CLUSTER_NAME)
	@echo "🚀 Applying Kubernetes manifests..."
	kubectl apply -f k8s/
    @echo "🔄 Forcing pod restart..."
    kubectl rollout restart deployment/http-server
	@echo "⏳ Waiting for deployment rollout..."
	kubectl rollout status deployment/http-server --timeout=60s
	@echo "✅ Server successfully deployed! Try: curl -v http://localhost/hello"
