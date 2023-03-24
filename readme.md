Service Description

This service is designed to receive requests, generate tasks from them, and send them to be executed by a third party. 
The standard Golang net/http library was used for implementation. The service runs on port 9090 and implements two methods:

    POST /task - for receiving and processing tasks.
    GET /task/taskid - for returning the task execution status.

Add support Redis and Postgress.

Add config
  Before run service you need check config/config.yaml file, and make sure all parameters are filled in correctly

postgres: 
  dsn: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
redis:
  addr: localhost:6379
  password: ''
  db: 0        
server: 
  port: "9090"
  storage: "postgres"

 storage - is nidicate where tasks will be saved

"postgres" - in Postgress Database
"redis" - in Redis Database
"memory" - in map (key\value datastruct golang)

If selected "postgres" you need create table :

CREATE TABLE tasks (
    id varchar(255) PRIMARY KEY,
    method varchar(255),
    url varchar(255),
    headers jsonb,
    status varchar(255),
    http_status_code integer,
    headers_array text[],
    length integer
);






