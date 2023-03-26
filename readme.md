README

Service Description:

The service is designed to receive requests, generate tasks from them, and send them to be executed by a third party. It implements two methods: POST /task and GET /task/taskid. The standard Golang net/http library was used for implementation. The service runs on port 9090.

Support for Redis and Postgres was added. Before running the service, please ensure that all parameters are filled in correctly in the config/config.yaml file.

The config file contains the following parameters:

- postgres:
    - dsn: "<database_connection_string>"
- redis:
    - addr: "<redis_server_address>'
    - password: '<redis_password>'
    - db: <redis_db_index>
- server:
    - port: "<port_number>"
    - storage: "<postgres/redis/memory>"

The `storage` parameter indicates where tasks will be saved: 

- If `postgres` is selected, a table named `tasks` should be created with the following schema:

    ```
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
    ```

- If `redis` is selected, tasks will be saved in a Redis database.

- If `memory` is selected, tasks will be stored in a map (key/value data structure in Golang).

**Note:**

- Please ensure that the PostgreSQL and Redis servers are running and accessible before running the service.
- Make sure to update the config.yaml file with the correct database connection string and Redis server credentials if needed.

TODO
  - Goose migrations
  - Prometeus metrics