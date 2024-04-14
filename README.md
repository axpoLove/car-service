# Car API Service

## Features

- get a list of cars by filter with pagination
- add a new cars by reg nums
- update a car
- delete a car

## How to run

1. Clone the repository:
   ```bash
   git clone https://github.com/some_name/car-api-service.git
   cd car-api-service
2. rename .env.example to .env and set all variables
3. run the following commands to run service
    ```bash
    make docker
    make run-car-info-api
    make run-car-api

## Documentation
To open the swagger documentation follow the links:
http://localhost:{CAR_SERVER_PORT}/swagger/
http://localhost:{CAR_INFO_SERVER_PORT}/swagger/
