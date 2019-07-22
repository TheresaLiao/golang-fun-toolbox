IMAGE_NAME=golang-gin:v1.0
CONAINER_NAME=golang-gin

docker stop $CONAINER_NAME
docker rm $CONAINER_NAME

docker build -t $IMAGE_NAME .
docker run -ti -d \
-p 80:80 \
--name $CONAINER_NAME \
$IMAGE_NAME 
