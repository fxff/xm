### Prerequisites
- docker-compose
- go 1.19

### Run
Just `make` 

### How to use
First you must authentificate:
```bash
curl -vvv -X GET localhost:8080/auth/ -d '{"user":"1","password":"1"}'
```
Pass same username and password to gain access.

Then you are ready to create companies:
```bash
curl -vvv -X POST localhost:8080/company/ -d '{"name":"home", "registered":false, "employees": 3, "type":"non-profit"}' --cookie "token=..."
# you will see something like {"ID":"cd2b510a-944f-408e-a16e-e2631e2bd32d"}
# use it to get information about company:
curl -vvv -X GET localhost:8080/company/cd2b510a-944f-408e-a16e-e2631e2bd32d --cookie "token=..."
```

Available types are:
- "corporation"
- "non-profit"
- "cooperative"
- "sole-proprietorship"
