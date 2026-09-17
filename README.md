# docker-compose-actions-workflow
[![Actions Status](https://github.com/peter-evans/docker-compose-actions-workflow/workflows/docker-compose-actions-workflow/badge.svg)](https://github.com/peter-evans/docker-compose-actions-workflow/actions)

This is a GitHub Actions workflow example to demonstrate building and testing a multi-container stack using `docker compose`.

The sample app is a small Go web service backed by Redis (migrated from the original Python/Flask sample based on the [Get started with Docker Compose](https://docs.docker.com/compose/gettingstarted/) documentation).

## Running locally

```sh
docker compose up -d --build
go test ./...
```

The integration tests in `main_test.go` require the stack to be running (via `docker compose up -d`) and hit `http://localhost:5000/` (override with the `APP_URL` env var).

## GitHub Actions Workflow

**push.yml**
```yml
name: Docker Compose Actions Workflow
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Build the stack
        run: docker compose up -d --build
      - name: Test
        run: go test ./...
```

You can browse a run for this example [here](https://github.com/peter-evans/docker-compose-actions-workflow/actions/workflows/push.yml).

For more about testing containers before release see [Smoke Testing](https://github.com/peter-evans/smoke-testing).

## License

[MIT](LICENSE)
