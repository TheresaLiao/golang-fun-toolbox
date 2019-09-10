imagename=u2bdownload:temp
containername=u2bdownload

docker stop $containername
docker rm   $containername
docker build -t $imagename .

docker run -d -it \
--name $containername \
$imagename
