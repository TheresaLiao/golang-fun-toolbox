
IMAGE_NAME=192.168.8.25:31115/golang-gin
TAG_NAME=v4.1
CONAINER_NAME=golang-gin

kubectl delete -f dnn-k8s-api.yaml

docker stop $CONAINER_NAME
docker rm $CONAINER_NAME

docker build -t $IMAGE_NAME:$TAG_NAME .
docker push $IMAGE_NAME:$TAG_NAME

sed -e "s|test|$TAG_NAME|g" dnn-k8s-api-tmp.yaml >> dnn-k8s-api.yaml


kubectl create -f dnn-k8s-api.yaml
