## Knowledge
# image is a blueprint
# container is a running instance of the image

## Build
# first layer, initial layer, or parent image
# pull golang as this layer
# when we create image this will pull golang first
# its going to pull it either from docker hub repository or our computer if we have already downloaded it
# Alpine is chosen for its small footprint compared to Ubuntu
FROM golang:1.18-alpine

# working directory
# after this instruction then if we specify path / if when we run image in the future it will run in "/app" directory as root (.)
WORKDIR /app

# additional layer - dependencies layer
# install all go dependencies
# RUN will run the command as the image is being built / at the build time
# positioned before copy all file because we need this layer to be cached for faster build time
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# additional layer - source code layer
# copy all files from source (.)(current machine) to destination (. <relative to working directory which is "/app">)(container image)
# use "." as source since Dockerfile is on the same level with the sourcecode
# use "/app" as destination to make sure that our source code file not conflicting with other files on the same images
COPY . .

# build 
RUN go build -o /oceana-app

# EXPOSE will tell the container which port should the container export
# We need to expose the port because if we run the applicaiton inside the containers then the application will run on the containers port not the actual machine port
EXPOSE 8080

# run the application in the container
# CMD will run the command on container
CMD ["/oceana-app"]
