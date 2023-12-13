# How to run with curl

```sh

➜ curl http://localhost:8080/users
[{"id":"1","name":"bob"}]

➜ curl http://localhost:8080/users/1
{"id":"1","name":"bob"}

➜ curl -X POST -H 'content-type: application/json' --data '{"id": "2", "name": "karen"}' http://localhost:8080/users
{"id":"2","name":"karen"}

➜ curl http://localhost:8080/users
[{"id":"1","name":"bob"},{"id":"2","name":"karen"}]

➜ curl http://localhost:8080/users/2
{"id":"2","name":"karen"}

```
