# Patika E-commerce API

This repository contains E-commerce API written by Go. It is a RESTful API.
Authentication is done using JWT. It is a JSON Web Token.
System has access and refresh tokens. 
The tokens validation time can be configured in the config file.

App has three different roles which are:
Admin, User and Anonymous.

Admin user can create, update, delete and get products and categories.
And also can create bulk categories via upload csv file.

Anonymous user can list and search products via pagination.

Authenticated user can;
 - create cart, add to cart and remove from cart.
 - list cart items and his/her orders.
 - When they submit the cart, the cart is checked for validity and if it is valid,
    the order is completed and the cart is paid.
 - If the order submited date is older than 14 days, the order can be canceled.

## Using Tools
 - Gin
 - Gorm
 - Postgres
 - JWT
 - net/http
 - viper
 - swagger
 - zap (logger)

## Clone the project
```
$ git clone https://github.com/Picus-Security-Golang-Bootcamp/bitirme-projesi-eibrahimarisoy.git
$ cd bitirme-projesi-eibrahimarisoy.git
```

## Dependencies used in the project
 - Go 1.18
 - Postgres 14.2

## Configuration
App uses viper to read configuration file.
You can edit file `config.yaml` in the project pkg/config directory.
There is a sample file `config.yaml.example` in the project pkg/config directory.

## Run the tests
```
$ go test ./...
```
## Run the project
```
$ go run main.go
```

## Run the project with swagger
Installing swagger please follow the instructions on the link below.

https://goswagger.io/install.html

```
$ swagger serve ./docs/patika-ecommerce.yml
```

## Routes
Default **Patika E-commerce API** routes are listed below. 

| METHOD  | ROUTE                           | DETAILS                                         |
|---------|---------------------------------|-------------------------------------------------|
| POST    | /api/v1/register                | user register endpoint                          |
| POST    | /api/v1/login                   | user login endpoint                             |
| POST    | /api/v1/refresh                 | refresh token endpoint                          |
| POST    | /api/v1/categories              | category create endpoint (admin)                |
| GET     | /api/v1/categories              | category list endpoint                          |
| GET     | /api/v1/categories/:id          | category detail endpoint                        |
| PUT     | /api/v1/categories/:id          | category update endpoint (admin)                |
| DELETE  | /api/v1/categories/:id          | category delete endpoint (admin)                |
| POST    | /api/v1/categories/bulk-upload  | category bulk upload endpoint (admin)           |
| GET     | /api/v1/products                | product list endpoint                           |
| GET     | /api/v1/products/:id            | product detail endpoint                         |
| POST    | /api/v1/products                | product create endpoint (admin)                 |
| PUT     | /api/v1/products/:id            | product update endpoint (admin)                 |
| DELETE  | /api/v1/products/:id            | product delete endpoint (admin)                 |
| POST    | /api/v1/cart                    | get or create cart endpoint (authenticated user)|        
| POST    | /api/v1/cart/add                | add to cart endpoint (authenticated user)       |
| GET     | /api/v1/cart/items              | list cart items endpoint (authenticated user)   |
| PUT     | /api/v1/cart/items/:id          | update cart item endpoint (authenticated user)  |
| DELETE  | /api/v1/cart/items/:id          | delete cart item endpoint (authenticated user)  |
| POST    | /api/v1/orders                  | complete order endpoint (authenticated user)    |
| GET     | /api/v1/orders                  | list orders endpoint (authenticated user)       |
| PUT     | /api/v1/orders/:id              | cancel order endpoint (authenticated user)      |
| GET     | /api/v1/healthz                 | application health check endpoint               |
| GET     | /api/v1/readyz                  | application readiness check endpoint            |

## Contact

If you want to contact me you can reach me at <eibrahimarisoy@gmail.com>.