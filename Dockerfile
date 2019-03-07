
#REF: https://docs.docker.com/engine/reference/builder/
#*********************************
# STEP 1 build executable binary
#*********************************
# use golang:alpine image to build the image, which contains all the 
#correct build tools and libraries. "As builder" its a container name for reference later on.
FROM golang:alpine AS builder
#Add Maintanier info
LABEL maintainer="Gangadhar Pathipaka"

#ADD . /go/src/go-docker/consignment-service
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

#set the current working directory inside the container
WORKDIR /go/src/consignment-service-mgo

#copy everything from current directory to the WORKDIR
COPY . .

#download all dependencies 
RUN go get -d -v ./...

#Build the Binary with flags which will allow us to 
#RUN this binary on Alpine.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/consignment-service-mgo

#*********************************
# STEP 2 build a small image 
#*********************************
#telling Docker to start another build process with this image.
FROM scratch

#Copy our static executable.
#Here, instread of copying the binary from our host machine 
#we pull the binary from the container named "builder", within 
#this build context. this reaches into previous image and finds the binary we build
#pulls it into this container. 
COPY --from=builder /go/bin/consignment-service-mgo /go/bin/consignment-service-mgo

EXPOSE 9000

#Run the consignment-service Binary
#run the binary buil in a separate container, with all of correct dependencies and 
#run this libraries.
ENTRYPOINT [ "/go/bin/consignment-service-mgo" ]