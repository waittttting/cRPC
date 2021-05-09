GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o cx-control-center main.go
docker build -t registry.local.com/cx-control-center .
docker push registry.local.com/cx-control-center
kubectl delete deployment cx-control-center -n cx-rpc-base
kubectl apply -f cx-control-center-deployment.yaml
docker rmi registry.local.com/cx-control-center
rm cx-control-center
