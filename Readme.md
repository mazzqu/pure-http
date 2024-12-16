# How to Execute locally

```sh
go run main.go
```
## Respons code
```
200 OK 
404 Not found
```

## Get list's of data


GET http://localhost:8080/users/


- For specsific UserId use :

GET http://localhost:8080/users/{id}



## Post a data

POST http://localhost:8080/users/

@body_parameters:
```json
{"id":"2","name":"who is you",}
```

## Update a data

PUT http://localhost:8080/users/{id}

@body_parameters:
```json
{"id":"2","name":"who is you",}
```

noted : When Update data contents can be duplicated.

## Delete a data

DELETE http://localhost:8080/users/{id}