# Simple Microservice

## Authors

- [Shay Jacoby](https://github.com/shayja)

## Description

A simple microservice project using Go, Gin and PostgreSQL.

- Go : https://go.dev/doc/
- Gin : https://gin-gonic.com/docs/
- PostgreSQL : https://www.postgresql.org/docs/

If running the app locally, ensure that PostgreSQL is installed on your machine. Alternatively, if you run the app using Docker Compose, everything will be automatically set up within the container.

Database:

1. Create a new PostgreSQL database named "shop".
2. Add a new db user called "appuser" and assign a login password.
3. Execute the SQL script located in the /migrations directory of the project on the "shop" database.
4. Update your database credentials in the .env.local file, then rename the file to .env. Do not move this file from root folder.
5. Adjust the configuration values to match the details of your "appuser" and the database root admin user.

# Database settings:

DB_HOST="localhost"
DB_USER="appuser"
DB_PASSWORD="<YOUR_PASSWORD>"
DB_NAME="shop"
DB_PORT=5432

configure the Postgres admin user credentials:
PGADMIN_DEFAULT_EMAIL="your@admmin.email.here"
PGADMIN_DEFAULT_PASSWORD="<<PGADMIN_ADMIN_PASSWORD>>"

To start the app, open Terminal:
go run ./cmd/main.go

To start using docker compose:
docker compose up --build

To stop the container:
docker-compose down --remove-orphans --volumes

## App endpoints:

**GET**
/api/v1/order

Get all products with paging

example req:
curl --location 'http://localhost:8080/api/v1/product?page=1' \
--data ''

example call:

**POST**
curl --location '/api/v1/order' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzUwNzM5NjQsImlhdCI6MTczNDk4NzU2NCwic3ViIjoiNDUxZmE4MTctNDFmNC00MGNmLThkYzItYzlmMjJhYTk4YTRmIn0.DNQ6qOa_b0Y8eLOEB5i0Es-4kBYWsEVVaTNWgPn-VNQ' \
--data '{
"order_details": [
{
"product_id": "063d0ff7-e17e-4957-8d92-a988caeda8a1",
"quantity": 1,
"unit_price": 101.00,
"total_price": 102.00
}
],
"status":1,
"total_price": 102,
"user_id": "451fa817-41f4-40cf-8dc2-c9f22aa98a4f"
}'
