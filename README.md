# User management (CRUD) app REST API 

The aim of this project is to create a scalable restful API that handles user management,
including user creation, retrieval, update, and deletion.

## Built with

- [gin](https://github.com/gin-gonic/gin) - The web framework used
- [viper](https://github.com/spf13/viper) - Configuration management
- [di](https://github.com/sarulabs/di) - Dependency injection framework
- [zap](https://pkg.go.dev/go.uber.org/zap) - Logging

## Notes, design decisions, assumptions made

- `domain.ID` could be also a part of the pkg to be used across different domains or layers, however here its use is tightly coupled with the user domain model
- `hashPassword()` function defined close to its usage
- addresses types are being checked against uniqueness to prevent duplicate entry errors
- it would be nice to override validation errors to use some more end-api-user friendly communication 
- there are different approaches possible and it really depends on the use case how to build the API contract or tackle the update action.  
It would be the best to clarify business requirements yet I need to do some assumptions in order to not waste time on async communication
- there are different approaches in terms of validation responsibility. In my solution I want to provide some fundamental validation and sanity checks for the input data at the infrastructure layer 
and in the application layer there is more business-related validation 
- I design the business based on the principle of trying and failing rather than checking if everything is ready for successful execution, if possible. This approach helps to minimize the number of calls
- There are some domain specific errors defined as global vars but those rather technical remain just non defined, yet informing about the root cause
- `ExecContext()` already uses prepared statements to prevent SQL injection
- instead of leaking error message at the output I log 500 errors in HTTP containers. Application layer log errors at debug level
- I could consider different approach about responses. Some APIs return created objects, I respond with UUID only. In order to get the real values from DB it'd require additional call. 
I don't think it's necessary, but it all depends on the business requirements
- I do not remove anything from database, that is a bad practice. I use soft deletes instead
- For simplicity I use the same database for service and integration test, however it would deserve dedicated db and seeds 
- Benchmark test is a draft
- It's worth doing memory dump and review allocations. Maybe it would be good to decrease number of allocations and sacrifice code maintainability to achieve performance
- I could try make GetByUUID() and GetAddresses() concurrent when I fetch user from DB in application's service
- I think "get user" and "get multiple" users deserve caching

## Project structure

This project is organized in a way that keeps different parts of the code separate and easy to manage, making it simpler to maintain and scale as needed for future.


```bash
/cmd
└── /server
    └── main.go

/internal
├── /config
│   └── config.go             # app's configuration, loading settings
│
├── /domain
│   └── /user
│   │   ├── model.go          # User model/entity and related domain logic
│   │   ├── repo.go           # repository interfaces, abstracting data access patterns
│   └── (...)
│
├── /application
│   └── /service
|      └── user.go            # service interface and implementation, business and data flow controller
|      └── (...)
│
└── /infrastructure
    ├── /container
    │   └── container.go      # dependency injection
    ├── /httpserver
    │   ├── server.go         # contains server and routing
    │   └── /handlers         # HTTP handlers per each resource interacting with app services
    │       └── http_user.go
    │       └──(...)    
    └── /database             # data access patterns implementations     
        └── /mysql
            └── repo_user.go
            └── /migrations   # schema migrations
            └── /entity       # db responses entities
            └── (...)

/pkg                          # sharable utils
└── /logger
    └── logger.go             

```

## Usage

```bash
git clone https://github.com/wojciechpawlinow/usermanagement.git && cd usermanagement

make build # build the docker image

make run # create containers
```

#### Run locally

```bash
docker-compose up mysql -d

migrate -path internal/infrastructure/database/mysql/migrations -database "mysql://user:pass@tcp(localhost:3306)/users" up
```

1. Create file _config.yaml_ with the following content
```yaml
DB_READ_HOST: localhost
DB_WRITE_HOST: localhost
```
and run 
```bash
go run cmd/server/main.go
```

2. Simply pass those values as ENVs
```bash
DB_READ_HOST=localhost DB_WRITE_HOST=localhost go run cmd/server/main.go
```

### Example cURLs

See [API docs](docs/api.md)

## Testing

Run unit tests
```bash
make test
```

Run benchmark
```bash
DB_READ_HOST=localhost DB_WRITE_HOST=localhost go test -run=^$ -bench=BenchmarkCreateUser ./tests
```

HTTP controller and application's service are covered with units.

There is also a sequence of calls being called against the docker's database (integration tests)

## Lint

```bash
make lint
```

## Database migrations

All migrations apply automatically while creating containers. 

If you want to drop and recreate schema, then
```bash
make reset-db
```

To drop data
```bash
mysql -h 127.0.0.1 -P 3306 -u user -ppass users
drop table addresses;
drop table users;
drop table schema_migrations;
```

Connect to docker MySQL 
```bash
mysql -h 127.0.0.1 -P 3306 -u user -ppass users

mysql> show tables;
+-------------------+
| Tables_in_users   |
+-------------------+
| addresses         |
| schema_migrations |
| users             |
+-------------------+
3 rows in set (0,00 sec)

mysql> describe addresses; describe users;
+-------------+--------------+------+-----+-------------------+-------------------+
| Field       | Type         | Null | Key | Default           | Extra             |
+-------------+--------------+------+-----+-------------------+-------------------+
| id          | bigint       | NO   | PRI | NULL              | auto_increment    |
| user_id     | bigint       | NO   | MUL | NULL              |                   |
| street      | varchar(255) | NO   |     | NULL              |                   |
| city        | varchar(100) | NO   |     | NULL              |                   |
| state       | varchar(100) | YES  |     | NULL              |                   |
| postal_code | varchar(20)  | NO   |     | NULL              |                   |
| country     | varchar(100) | YES  |     | NULL              |                   |
| created_at  | datetime     | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
| updated_at  | datetime     | YES  |     | NULL              |                   |
| deleted_at  | datetime     | YES  |     | NULL              |                   |
| type        | varchar(255) | NO   |     | NULL              |                   |
+-------------+--------------+------+-----+-------------------+-------------------+
11 rows in set (0,01 sec)

+--------------+--------------+------+-----+-------------------+-------------------+
| Field        | Type         | Null | Key | Default           | Extra             |
+--------------+--------------+------+-----+-------------------+-------------------+
| id           | bigint       | NO   | PRI | NULL              | auto_increment    |
| uuid         | char(36)     | NO   | UNI | NULL              |                   |
| email        | varchar(255) | NO   | UNI | NULL              |                   |
| password     | text         | NO   |     | NULL              |                   |
| created_at   | datetime     | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
| updated_at   | datetime     | YES  |     | NULL              |                   |
| deleted_at   | datetime     | YES  |     | NULL              |                   |
| first_name   | varchar(255) | NO   |     | NULL              |                   |
| last_name    | varchar(255) | NO   |     | NULL              |                   |
| phone_number | varchar(20)  | YES  |     | NULL              |                   |
+--------------+--------------+------+-----+-------------------+-------------------+
10 rows in set (0,01 sec)

```

If you want to provide changes by manually modyfing migration files once there already were copied to the container during the image creation process you must execute migration command from your host:  
```bash
migrate -path internal/infrastructure/database/mysql/migrations -database "mysql://user:pass@tcp(localhost:3306)/users" drop 1
migrate -path internal/infrastructure/database/mysql/migrations -database "mysql://user:pass@tcp(localhost:3306)/users" up
```

## Security

In terms of authentication I'd provide an integration with Auth0 and add a middleware that checks for JWT tokens to retrieve identity.  

For internal communication between microservices or external services I'll recommend M2M tokens. In terms of authorization we can use Polar language for defining rules of access and Oso framework to enable authorization in our service.
However, maybe https://github.com/casbin is a good alternative as well.