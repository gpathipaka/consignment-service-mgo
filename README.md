#build#
#protobuf 
protoc --go_out=plugins=micro:. proto/consignment-service/consignment.proto

#Docker Build
docker build -t cservice .

#docker Run
docker run -p 50051:50051 -e MICRO_SERVICE_ADDRESS="50051 -e MICRO_REGISTRY=mdns cservice