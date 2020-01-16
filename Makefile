
all:
	go fmt ./...
	go test ./...
	go mod tidy

minikube: all
	eval $$(minikube docker-env) && docker build -t poc-go-argo-pipeline .
	kubectl config use-context minikube
	kubectl delete -f minikube.yaml || echo ignore error
	kubectl wait --for=delete pod/example || echo ignore error
	kubectl apply -f minikube.yaml
