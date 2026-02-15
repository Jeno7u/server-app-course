## This project is result of test 1 in server app development course.

## Project structure

    cmd/app/main.go - entrypoint, router creation
    internal/
        models/
            feedback.go
            user.go
        src/
            index.html
        validation/ - custom validator		
            validation.go
	
---
## Run local

```
go run cmd/app/main.go
```
---
## Endpoints

`GET /` - returns index.html

`POST /calculate` - takes payload in format below and returns sum
```
{
	"num1": int
	"num2": int
}
```

`GET /users` - returns JSON below
```
{
	"id": 1
	"name": "Mironov Boris"
}
```

`POST /user` - takes payload in format below and returns same payload with extra field "is_adult". "is_adult" is true if age >= 18 else false
```
{
	"name": str
	"age": int
}
```

`POST /feedback` - takes payload in format below, validates it, adds it to list of feedbacks and returns {"message": "Feedback received. Thank you, {payload.name}"}

```
{
	"name": str
	"message": str
}
```

