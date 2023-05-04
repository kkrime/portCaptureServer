This repository is the solution, [this](https://github.com/kkrime/commandLinePoker/blob/main/developer-assignment-BA-go.md) is the problem
# portCaptureServer

**Please note**: I have not worked on a project that used Domain Driven Design before. </br>I familiarized myself with DDD and its core concepts as apart from this
assignment, and I tried to follow the DDD principles as best as I understood them.</br>

There are two golang programs in this repository:

## 1. portCaptureServerTranslator
This is a translation layer (from REST to gRPC).

**This can be thought of as outside the main scope of this assignment and is something auxiliary.**

## 2. portCaptureServer
This is the main program which saves the ports to the database<br>

### Design

#### API interface
For receiving the ports data **gRPC streams** are used, reasons:
- for microservices communications, gRPC is very efficient and is widely supported
- gRPC streams allow the service to control the flow of data and how much the service consumes at any one time. This is helpful for limiting resources.

**NOTE:** The default message size for gRPC is 4MB. I feel like this is a good enough size for most scenarios and most hardware portCaptureServer will be run on, so I did not include an option to set it in the config file.

#### Worker Threads
On start up a number of go routines are spawned, these are the worker threads that write the incoming port data to the database.
</br>
The number of worker threads that are used is defined in the config files located in: `portCaptureServer/config/` 
and can be adjusted to suite different hardware capabilities.
</br></br>
**Please note:** `portCaptureServer/app/service/*` is where the bulk of the functionality is coded.

#### Database
For the database I decided to go with Postgresql.</br>
The schema is located in `portCaptureServer/db/schema.sql`.
</br></br>
The port data inside `ports.json` was a little confusing:
```
"AEAJM": {
    "name": "Ajman",
    "city": "Ajman",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.5136433,
      25.4052165
    ],
    "province": "Ajman",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAJM"
    ],
    "code": "52000"
  }
```

As you can see from the above, each object is referenced by its `unloc`, but each object also has list of `unlocs` which contains the same `unloc` that the object is referenced with.</br>
This would imply that it is possible that a port can have multiple `unlocs`. 
</br>I made that assumption and also that there is a `"main unloc"` the one the object is referenced with, and I called this the `primary_unloc` in the schema/code. 
#### Saving to the database
For security and audit reasons, no data in the database is ever truly deleted, instead it is just marked as deleted.</br>

With regard to saving the ports data to the database, the service works in an all or nothing way; either all the data in a request is written to the database, or none of it is.
This means that if there is a single error in any of the ports, the whole file will have to be sent over again, there is no 'partial success'.</br>
The reason for this is to avoid database bloat and having multiple duplicate identical records (identical apart from the `deleted_at` field).

## Getting Started
Run the following from the root directory
1. First you need to apply the database schema (note you need psql installed for this to run):</br>
`sudo ./setupDockerDB.sh`
2. Build the docker images:</br>
`sudo docker-compose build`
3. Start all the services plus the database:</br>
`sudo docker-compose up` (the database needs to initialize before all the other services start, so this might take a minute)


## Testing
#### Integration testing
You will need to have `python3` and `pytest` installed:</br>
From `integrationTests` run `sudo python3 -m pytest -s -v`
#### Unit testing
I made sure to unit test the most critical parts of this service:</br>
From `portCaptureServer/app/service/` run `go test .`
