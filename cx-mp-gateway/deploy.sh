GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o cx-mp-gateway main.go
docker build -t registry.local.com/cx-mp-gateway .
docker push registry.local.com/cx-mp-gateway
kubectl delete deployment cx-mp-gateway-deployment -n cx-rpc-business
kubectl apply -f cx-mp-gateway-deployment.yaml
docker rmi registry.local.com/cx-mp-gateway
rm cx-mp-gateway
