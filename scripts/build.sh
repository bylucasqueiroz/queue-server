eval $(minikube docker-env)  # Configure the terminal to use Docker with Minikube
docker build -f cmd/server/Dockerfile -t server:latest .
docker build -f cmd/client/Dockerfile -t client:latest .