# Go Prometheus Client Example

This is an example of a web app that queries a Prometheus server. Originally intended as a proof-of-concept for querying 
the Prometheus server in the using [Rancher](https://github.com/rancher/rancher)'s cluster monitoring, this can also be
applied to other Prometheus instances that optionally have basic authentication.

The web app uses [Echo](https://echo.labstack.com/).

## Getting Started

For development, using dotenv is recommended for configuration. In production, this can be replaced with docker-compose
`environment` or kubernetes config map.

```shell script
cp .env.example .env
```

- `PROMETHEUS_URL`: Prometheus URL
- `PROMETHEUS_AUTH_USERNAME`: Optional basic auth username
- `PROMETHEUS_AUTH_PASSWORD`: Optional basic auth password

The project comes with basic support for docker container development and deployment. This is optional.

```shell script
cd docker
cp .env.example .env
docker-compose up -d
```

[Reflex](https://github.com/cespare/reflex) will automatically reload the server if there are changes in `.go` or `.html` files.

## Rancher Cluster Monitoring 

To query the Prometheus instance deployed with Rancher's cluster monitoring:

- Deploy monitoring on a cluster
- Create a user role having the following grants on kubernetes resources: Create, Get, List, Watch on Resources `services,services/proxy,endpoints (Custom)`, API Groups `*`
- Grant the created role on a user
- Generate an API for the user
- Obtain the Prometheus URL. This can be acquired by going to the Prometheus web interface. Example URL: http://rancher.local/k8s/clusters/c-pw2r2/api/v1/namespaces/cattle-prometheus/services/http:access-prometheus:80/proxy

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
