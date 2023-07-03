# rings-social-backend

This is the backend code for [rings.social](https://rings.social) a content-voting platform that is Reddit-API compatible.  

## Requirements

- Go
- `make`
- Docker
- `docker-compose`
- Auth0 Application (auth can't be disabled for now)

## Getting started

```
docker-compose up -d
make build

# Config
export DATABASE_URL=postgresql://ring:ring@localhost:5432/ring
export AUTH0_DOMAIN=your-domain.auth0.com
export AUTH0_CLIENT_ID=xyz
read -s -r AUTH0_CLIENT_SECRET # type the secret and press ENTER
export AUTH0_CLIENT_SECRET

./build/rings-backend
```

Congrats! The backend should be up and running on the displayed address.
To listen on `0.0.0.0` just pass the `-l` argument (for example `-l 0.0.0.0:8080`)

## Testing it

Choose one of the [routes](./pkg/routes.go) or
modify a Reddit application to point to your backend.  

Alternatively, use the [Rings frontend](https://github.com/rings-social/frontend) to 
have the full [rings.social](https://rings.social) experience.

## Reddit compatible API

We're planning to have a Reddit compatibility layer to allow the existing apps (e.g: Sync, RIF, Apollo, ...)
to effortlessly migrate to Rings.


## Contributions 

Contributions are welcome. Help us shape the future by sending your PRs!
