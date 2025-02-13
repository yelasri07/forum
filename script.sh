# Names for our image & container
IMG=my_img
CTR=my_ctr


# Remove all Docker containers
docker rm -f $(docker ps -aq)

# Build Our image, (--no-cache)=> each build step will be executed without retrieving already stored data.
docker build --no-cache -t $IMG .

# Run the Docker container and link port 8080 on the container to port 8080 on your machine
# Run container -d => This option starts the container in detached mode (background)
docker container run -d -p 8080:8080 --name $CTR $IMG
