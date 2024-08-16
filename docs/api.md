## Example calls

### Create user 
```bash
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '
{
  "email": "test1@gmail.com",
  "password": "admin123",
  "first_name": "Test1",
  "last_name": "Test1",
  "phone_number": "111111111",
  "addresses": [
    {
      "type": 1,
      "street": "Test1",
      "city": "Test1",
      "state": "Test1",
      "postal_code": "11111",
      "country": "USA"
    },
   {
      "type": 2,
      "street": "Test1",
      "city": "Test1",
      "state": "Test1",
      "postal_code": "111111",
      "country": "USA"
    }
  ]
}'
```
Response:
```bash
{"uuid":"495e962a-51db-4d38-bfbe-048254022d9d"}
```

### Get user by identifier

```bash
curl -X GET http://localhost:8080/users/495e962a-51db-4d38-bfbe-048254022d9d
```
Response
```bash
{
  "id": "495e962a-51db-4d38-bfbe-048254022d9d",
  "email": "test1@gmail.com",
  "first_name": "Test1",
  "last_name": "Test1",
  "phone_number": "111111111",
  "addresses": [
    {
      "type": 1,
      "street": "Test1",
      "city": "Test1",
      "state": "Test1",
      "postal_code": "11111",
      "country": "USA"
    },
    {
      "type": 2,
      "street": "Test1",
      "city": "Test1",
      "state": "Test1",
      "postal_code": "111111",
      "country": "USA"
    }
  ]
}
```

### Get users with pagination

```bash
curl -X GET http://localhost:8080/users?size=3&page=1
```

### Update user
```bash
curl -X PUT http://localhost:8080/users/495e962a-51db-4d38-bfbe-048254022d9d -H "Content-Type: application/json" -d '{
  "first_name": "Test111111111",
  "last_name": "Test111111111",
  "phone_number": "123456789",
  "addresses": [
    {
      "type": 3,
      "street": "Test111111111",
      "city": "Test111111111",
      "postal_code": "111111"
    }
  ]
}'
```
Response
```bash
"ok"
```

### Delete user
```bash
curl -X DELETE http://localhost:8080/users/495e962a-51db-4d38-bfbe-048254022d9d
```
Response 
```bash
"ok"
```
