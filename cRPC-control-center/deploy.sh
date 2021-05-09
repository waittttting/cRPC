docker build -t cx-control-center .
docker tag cx-control-center registry.local.com/cx-control-center
docker push registry.local.com/cx-control-center
kubectl delete deployment cx-contol-center-deployment -n cx-rpc-base
kubectl apply -f cx-control-center-deployment.yaml
docker rmi cx-control-center
docker rmi registry.local.com/cx-control-center
