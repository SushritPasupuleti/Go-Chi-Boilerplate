# Go-Chi Boilerplate

Creating boilerplate for Go-chi with some good defaults.

> [!NOTE]  
> This repository is updated infrequently, as I backport changes from my other projects into this. However, at any given point, the repo should be considered a good starting point.

## Features

- [x] Fully documented codebase with GoDoc.
- [x] Logging with [zerolog](https://https://github.com/rs/zerolog)
- [x] Routing with [go-chi](https://go-chi.io/)
- [x] OpenAPI with [go-swagger](https://github.com/swaggo/swag)
- [x] Input Validation with [go-playground/validator](https://github.com/go-playground/validator)
- [x] Sane HTTP Security Headers with [secure](https://github.com/unrolled/secure)
- [x] Custom Redis Cache Middleware with [go-redis](https://github.com/redis/go-redis)
  - [ ] Optional: Memcached implementation
- [x] OAuth 2.0 client.
  - [x] Password hashing with [bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt)
  - [x] Token Grant
  - [x] Token Validation + RBAC
  - [x] Token Refresh
  - [x] Token Revoke
- [x] JWT authentication.

## Setup

Run `make` to see all available commands.

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

## Notes on Design Considerations

- JWTIDs were used, but for the `refresh token` only. This is because the `refresh token` is persisted in the `redis` cache, and therefore needs to be revoked. `The access token` is not persisted, and therefore does not need to be revoked. This has the following benefits:

  - The `access token` doesn't need to be validated against the DB or cache, on each request. And instead the `refresh token` requires this only during a refresh.

  - This avoids too many DB/cache lookups, and therefore improves performance.

  - You can however, choose to use JWTIDs for both the `access token` and `refresh token`, if you want to prevent `replay attacks`.

## Notes on storage of JWTs

- The API returns both an `access token` and a `refresh token`, it is recommended that the `access token` is stored in memory, and the `refresh token` is stored in a cookie with the `secure` & `http-only` flags set.

- The `refresh token` is also persisted in the `redis` cache for validation and revocation.

- Persisting the `access token` in memory, means that the token is not persisted across browser restarts, and is therefore more secure.

- Your client should refresh the `access token` when it expires, using the `refresh token` stored in the cookie.

- In case of higher security requirements, you can follow one of the following patterns:

  - Use a `refresh token` with a short expiry time, and refresh the `refresh token` on every request.

  - Omit the `refresh token` entirely, and use a short lived `access token`, and prompt the user to login again when the `access token` expires.

- You may also consider using JWTIDs, to prevent replay attacks.
