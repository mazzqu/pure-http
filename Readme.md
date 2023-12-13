# How to run with curl

```sh

➜ curl http://localhost:8080/users
[{"id":"1","name":"bob"}]

➜ curl http://localhost:8080/users/1
{"id":"1","name":"bob"}

# No space between 'content-type: application/json'.
# replaced the single quotes with double quote and escaped the double quotes inside the curly braces with -
# a backslash and it has worked.
# ➜ curl -X POST -H 'content-type: application/json' --data '{"id": "2", "name": "lolen"}' http://localhost:8080/users

➜ curl -X POST -H "content-type:application/json" --data "{\"id\": \"2\", \"name\": \"lolen\"}" http://localhost:8080/users
{"id":"2","name":"lolen"}

➜ curl http://localhost:8080/users
[{"id":"1","name":"bob"},{"id":"2","name":"lolen"}]

➜ curl http://localhost:8080/users/2
{"id":"2","name":"lolen"}

```
