GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o cx-config-center main.go
docker build -t registry.local.com/cx-config-center .
docker push registry.local.com/cx-config-center
kubectl delete deployment cx-config-center -n cx-rpc-business
kubectl apply -f cx-config-center-deployment.yaml
docker rmi registry.local.com/cx-config-center
rm cx-config-center
