build:
	docker build -t consignment-service-mgo .
	
run:
	docker run --net="host" \
		-p 50052 \
		-e MICRO_SERVER_ADDRESS=:50052 \
		-e MICRO_REGISTRY=mdns \
		-e DISABLE_AUTH=true \
		consignment-service-mgo
		
run_mongo:
	docker run -d -p 27017:27017 mongo