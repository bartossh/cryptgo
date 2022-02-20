docker run -d -it --name docker-dev-go --mount type=bind,source="$(pwd)",target=/app -p 8085:8085 -p 3000:3000 golang:latest sleep infinity
