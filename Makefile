
all:
	go fmt ./...
	go test ./...
	go mod tidy

argo:
	kubectl create ns argo || echo ignore error
	kubectl apply -n argo -f https://raw.githubusercontent.com/argoproj/argo/v2.4.3/manifests/install.yaml

minikube: all argo
	eval $$(minikube docker-env) && docker build -t poc-go-argo-pipeline .
	kubectl config use-context minikube
	kubectl delete -f minikube.yaml || echo ignore error
	kubectl wait --for=delete pod/example || echo ignore error
	kubectl apply -f minikube.yaml
