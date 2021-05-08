GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o cx-mp-user main.go
docker build -t registry.local.com/cx-mp-user .
docker push registry.local.com/cx-mp-user
kubectl delete deployment cx-mp-user -n cx-rpc-business
kubectl apply -f cx-mp-user-deployment.yaml
docker rmi registry.local.com/cx-mp-user
