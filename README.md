[![Build Status](https://travis-ci.org/MYOB-Technology/shell-operator.svg?branch=master)](https://travis-ci.org/MYOB-Technology/shell-operator)

# Kubernetes Shell Operator

This operator is a generic operator that can watch any Kubernetes Object that you specify and will just execute a shell command in a subshell on any change to that Object.

The usecase of this operator is for Kubernetes Cluster Administrators to be able to automate any workflow in their cluster based on Kuberenetes Object change events without having to write the Kubernetes watch boilerplate everytime. **It is NOT intended as a way for any Kubernetes user to send an arbitrary shell command into the shell-operator pod.**

## Goals

- Centralise the development and testing of Kubernetes CRD specific controller workflows to a single place so we dont duplicate the code for every operator
- Provide a robust operator pattern in code that is easy to leverage
- Allow for quick prototyping of an operator in any language, or even bash scripts
- decouple the CRD watching from the operator actions to make it easy to test the actions without kubernetes or other dependencies

## How it works

To use this operator you create a new project with a Dockerfile that contains everything you need for your operator use. Just add the Shell Operator Docker image as a separate multistage `FROM` image and copy it into your Dockerfile. See [example/Dockerfile](example/Dockerfile) for an example.

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
    # Set env vars to be available in the shell
    # This way you can set environment specific items
    # as per a normal 12 factor app
    environment:
      DB_URL: "xxx"
```

### CRD input

The operator will expose environment variables into the shell environment when the script is run to allow the script to identify the namespace and name of the object that has changed. The values are:

* SHOP_OBJECT_NAMESPACE
* SHOP_OBJECT_NAME
* Any other variables you have listed in the `environment` key of the YAML config for that watch (see example above).

You can reference these environment variables in your script and use a Kubernetes API call or kubectl to get information on the object that has changed or is being reconciled.

## Development

The repo is enabled with docker and docker-compose. To run tests execute `docker-compose run test`.
