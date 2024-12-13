# How to run with curl

```sh

➜ curl http://localhost:8080/users/1
{"id":"1","name":"bob"}

➜ curl -X POST -H "content-type:application/json" --data "{\"id\": \"2\", \"name\": \"lolen\"}" http://localhost:8080/users
{"id":"2","name":"lolen"}

➜ curl http://localhost:8080/users
[{"id":"1","name":"bob"},{"id":"2","name":"lolen"}]

➜ curl http://localhost:8080/users/2
{"id":"2","name":"lolen"}

```

# How to run with thunder client

```http

GET -> http://localhost:8080/users

POST -> http://localhost:8080/users

## params @body for POST, e.g :
{["id":"1"],["name":"Ucup"]}

DELETE -> http://localhost:8080/users/{id}

```
