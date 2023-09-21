# Go-Chi Boilerplate

Creating boilerplate for Go-chi with some good defaults.

## Features

- [x] Logging with [zerolog](https://https://github.com/rs/zerolog)
- [x] Routing with [go-chi](https://go-chi.io/)
- [x] OpenAPI with [go-swagger](https://github.com/swaggo/swag)
- [x] Input validation with [go-playground/validator](https://github.com/go-playground/validator)
- [ ] OAuth 2.0 client.
  - [ ] Token Grant
  - [ ] Token Refresh
  - [ ] Token Revoke
- [x] JWT authentication.
- [ ] RBAC, ABAC with [casbin](https://pkg.go.dev/github.com/casbin/casbin/v2)

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
