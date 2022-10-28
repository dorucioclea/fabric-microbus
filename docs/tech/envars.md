# Environment Variables

The `Microbus` framework uses environment variables for various purposes:

* Initializing the connection to NATS
* Identifying the deployment environment (`PROD`, `LAB`, `LOCAL`)
* Designating a plane of communication

## NATS Connection

Before connecting to NATS, a microservice can't communicate with other microservices and therefore it can't reach the configurator microservice to fetch the values of its config properties. Connecting to NATS therefore must precede configuration which means that initializing the NATS connection itself can't be done using the standard configuration pattern. Instead, the [NATS connection is initialized using environment variables](./natsconnection.md): `MICROBUS_NATS`, `MICROBUS_NATS_USER`, `MICROBUS_NATS_PASSWORD` and `MICROBUS_NATS_TOKEN`.

## Deployment

`Microbus` recognizes three deployment environments:

* `PROD` for production deployments
* `LAB` for fully-functional non-production deployments such as dev integration, testing, staging, etc.
* `LOCAL` for when developing locally or when running applications inside a testing application

The deployment environment impacts certain aspects of the framework such as the log format and log verbosity level.

The deployment environment is set according to the value of the `MICROBUS_DEPLOYMENT` environment variable. If not specified, `PROD` is assumed, unless connecting to NATS on `localhost` in which case `LOCAL` is assumed.

## Plane of Communication

The plane of communication is a unique prefix set for all communications sent or received over NATS.
It is used to isolate communication among a group of microservices over a NATS cluster
that is shared with other microservices.

If not explicitly set via the `SetPlane` method of the `Connector`, the value is pulled from the `MICROBUS_PLANE` environment variable. The plane must include only alphanumeric characters and is case-sensitive.

Applications created with `application.NewTesting` set a random plane to eliminate the chance of collision when tests are executed in parallel in the same NATS cluster, e.g. using `go test ./... -p=8`.

This is an advanced feature and in most cases there is no need to customize the plane of communications.