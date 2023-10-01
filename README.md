# Go-Chi Boilerplate

Creating boilerplate for Go-chi with some good defaults.

## Features

- [x] Fully documented codebase with GoDoc.
- [x] Logging with [zerolog](https://https://github.com/rs/zerolog)
- [x] Routing with [go-chi](https://go-chi.io/)
- [x] OpenAPI with [go-swagger](https://github.com/swaggo/swag)
- [x] Input Validation with [go-playground/validator](https://github.com/go-playground/validator)
- [x] Custom Redis Cache Middleware with [go-redis](https://github.com/redis/go-redis)
  - [ ] Optional: Memcached implementation
- [ ] OAuth 2.0 client.
  - [x] Token Grant
  - [x] Token Validation + RBAC
  - [ ] Token Refresh
  - [ ] Token Revoke
- [x] JWT authentication.

## Setup

### Install dependencies

```bash
cd server
make packages_install
```

### Run

```bash
cd server
make run
```
