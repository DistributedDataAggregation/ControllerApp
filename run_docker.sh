sudo docker run --env-file .env \
    -p 3000:3000  \
    -v /home/data:/home/data \
    -it controller-image:latest