Service Description

This service is designed to receive requests, generate tasks from them, and send them to be executed by a third party. 
The standard Golang net/http library was used for implementation. The service runs on port 9090 and implements two methods:

    POST /task - for receiving and processing tasks.
    GET /task/taskid - for returning the task execution status.

Add config
Add support Redis and Postgress

