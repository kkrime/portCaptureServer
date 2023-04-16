# portCaptureServer

**Please note**: I have not worked on a project that used Domaon Driven Design before. </br>I familiarised myself with DDD and its core concepts as apart of this
assignment, and I tried to follow the DDD principles as best as I understood them.</br>

There are two golang programs in this repositroy:

## 1. portCaptureServerTranslator
This is a simple service that sends the ports data to the main service: portCaptureServer.</br>
The reason this is included is to provide a mean of communication (over gRPC) to portCaptureServer.</br>
</br>
**While it is assumed that the main service that writes to the database (portCaptureServer) could be run on modest hardware with limited RAM, 
this assumption does not hold for portCaptureServerTranslator, as portCaptureServerTranslator is assumed to be running on any machine and 
not even one what is neccessarily in the port domain** 

## 2. portCaptureServer
This is the main program which saves the ports to the database<br>

### Design

#### API interface
For recieving the ports data **gRPC streams** are used, reasons:
- for microservices communications gRPC is more efficent and is widly supported
- gRPC streams allows the service to control the flow of data and how much the service consumes at any one time. This is helpful for limiting resources.

#### Worker Threads
On start up a number of go routines are spawned, these are the worker threads that write the incoming port data to the database.
</br>
The number of worker threads that are used are defined in the config files located in: `portCaptureServer/config/` 
and can be adjusted to suite different hardware capabilities.
</br></br>
**Please note:** `portCaptureServer/app/service/*` is where the bulk of the functionality is coded.

#### Saving to the database
For secuirty and audit reaons, no data in the database is ever turely deleted, instead it is just marked as deleted.</br>

With regards to saving the ports data to the database, the service works in an all or nothing way; either all the data is written to the database, or non of it is.
This means that if there's a single error in any of the ports, the whole file will have to be sent over again, there is no 'partial success'.</br>
The reason for this is to avoid database bloat and having multiple duplicate identical records (identical apart from the `deleted_at` field).

## Getting Started
Run the following from the root directory
1. First you need to apply the database schema (note you need psql installed for this to run):</br>
`sudo ./setupDockerDB.sh`
2. Build the docker images:</br>
`sudo docker-compose build`
3. Start all the services plus the database:</br>
`sudo docker-compose up`


## Testing
#### Live testing
To run live testing, make sure docker-compose is running (from step 3 above), then from the root directory run:</br> `curl -X POST localhost:8080/v1/sendports -d @./testData/ports.json`
#### Unit testing
I made sure to unit test the most critical parts of this service:</br>
From `portCaptureServer/app/service/` run `go test .`
