[![Build Status](https://travis-ci.org/MYOB-Technology/shell-operator.svg?branch=master)](https://travis-ci.org/MYOB-Technology/shell-operator)

# Kubernetes Shell Operator

This operator is a Kubernetes controller framework that can watch any Kubernetes Object that you specify and will just execute a shell command in a subshell on any change to that Object.

The usecase of this operator is for Kubernetes Cluster Administrators to be able to automate any workflow in their cluster based on Kuberenetes Object change events without having to write the Kubernetes watch boilerplate everytime. **It is NOT intended as a way for any Kubernetes user to send an arbitrary shell command into the shell-operator pod.**

## Goals

- Centralise the development and testing of Kubernetes CRD specific controller workflows to a single place so we dont duplicate the code for every operator
- Provide a robust operator pattern in code that is easy to leverage
- Allow for quick prototyping of an operator in any language, or even bash scripts
- decouple the CRD watching from the operator actions to make it easy to test the actions without kubernetes or other dependencies

## How it works

The shell operator is a binary that you copy into your docker image. You can then copy or install any other depedencies and scripts you want to run for a change to Kubernetes objects. Once your docker image is deployed it will execute the shell operator binary and get the configuration you set to execute your shell scripts whenever the object you specify changes in the Kubernetes API Server.

To use this operator framework you create a new project with a Dockerfile that contains everything you need for your operator to use. Just add the Shell Operator Docker image as a separate multistage `FROM` image and copy it into your Dockerfile. See [example/Dockerfile](example/Dockerfile) for an example.

You can also set your config in a Kubernetes Config Map and volume it in to reduce how often you have to change your docker image.

### Configuration

The Shell Operator can be configured by creating a YAML configuration and copying it into your Docker image for it to pick up.

```yaml
---
# `boot` is optional
boot:
    # this command is run on boot and is useful for upserting your CRD creation object
    # or any other prep work to be done once when any new pod comes up.
  - command: kubectl apply -f /app/mycrd.yaml
    # optional key to set how long a process can be running in seconds before it is hard killed. The default is
    # 30 seconds. This makes sure the app will bootup within a reasonable time, or die of a timeout.
    timeout: 30
    # Set env vars to be available in the shell
    # This way you can set environment specific items
    # as per a normal 12 factor app
    environment:
      DB_URL: "xxx"

# `watch` is a required key and is an array of watches
watch:
    # The api group and version for the object you want to watch.
  - apiVersion: my.domain.io/v1beta1
    # The kind, these are the values that are in a yaml manifest representation of an object of this
    # type, so you can get the values from that.
    kind: MyObject
    # The command to run - the operator will execute this in a subshell with the default shell
    # it will pipe the Object being updated as a JSON object to stdin.
    command: python myscript.py
    # set how many workers the operator should use
    # there may be a need to serialise everything because of race conditions so this can be set to 1
    # or you might want multiple crds to be updated at the same time so
    # want to have 10 or more workers
    # This option is here to constrain the operator to the resources you want to use.
    concurrency: 1
    # optional key to set how long a process can be running in seconds before it is hard killed. The default is
    # 20 mins. This makes sure the concurrency is not exhausted and deadlocks the controller if processes
    # start freezing up.
    timeout: 1200
    # Set env vars to be available in the shell
    # This way you can set environment specific items
    # as per a normal 12 factor app
    environment:
      DB_URL: "xxx"
```

### CRD input

The operator will expose environment variables into the shell environment when the script is run to allow the script to identify the namespace and name of the object that has changed as well as the type. The values are:

* SHOP_OBJECT_NAMESPACE
* SHOP_OBJECT_NAME
* SHOP_API_VERSION
* SHOP_KIND
* Any other variables you have listed in the `environment` key of the YAML config for that watch (see example above).

You can reference these environment variables in your script and use a Kubernetes API call or kubectl to get information on the object that has changed or is being reconciled.

## Observability

Shell Operator has a metrics endpoint built in that you can point Prometheus at. This will allow you to see how many runs are occuring, which ones are failing and the average time in milliseconds on a per watch basis.

The endpoint is running on port 8080 at `/metrics` with a healthcheck at `/healthz` that you can use for liveness probes.

See [pkg/metrics/metrics.go](pkg/metrics/metrics.go) for the list of available metrics.

## Development

The repo is enabled with docker and docker-compose. To run tests execute `docker-compose run test`.
