### CREATE PVZ 

curl.exe -X POST "http://localhost:9000/pvz" -H "Content-Type: application/json" -d '{\"name\":\"пвз1\",\"address\":\"адрес1\",\"contact\":\"телефон1\"}'

curl.exe -X POST "http://localhost:9000/pvz" -H "Content-Type: application/json" -d '{\"name\":\"пвз2\",\"address\":\"адрес2\",\"contact\":\"телефон2\"}'

### UPDATE PVZ

curl.exe -X PUT http://localhost:9000/pvz/1 -H "Content-Type: application/json" -d '{\"name\":\"пвз11\",\"address\":\"адрес11\",\"contact\":\"телефон11\"}'

### LIST OF PVZ

curl.exe -X GET http://localhost:9000/pvz -H "Content-type: application/json"

### DELETE PVZ BY ID

curl.exe -X DELETE http://localhost:9000/pvz/2

### GET PVZ BY ID

curl.exe -X GET http://localhost:9000/pvz/1 -H "Content-type: application/json"