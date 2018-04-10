# Kubernetes Shell Operator

This operator is a generic operator that can watch any Kubernetes Object that you specify and will just execute a shell command in a subshell on any change to that Object.

The usecase of this operator is for Kubernetes Cluster Administrators to be able to automate any workflow in their cluster based on Kuberenetes Object change events without having to write the Kubernetes watch boilerplate everytime. **It is NOT intended as a way for any Kubernetes user to send an arbitrary shell command into the shell-operator pod.**

## Goals

- Centralise the development and testing of Kubernetes CRD specific controller workflows to a single place so we dont duplicate the code for every operator
- Provide a robust operator pattern in code that is easy to leverage
- Allow for quick prototyping of an operator in any language, or even bash scripts
- decouple the CRD watching from the operator actions to make it easy to test the actions without kubernetes or other dependencies

## How it works

To use this operator you create a new project with a Dockerfile like the following:

```dockerfile
FROM myobplatform/shell-operator:latest
# Install any dependencies you want including binaries
# The base image is build from alpine, so use `apk add -U ...`

# Add in your own code/scripts that you want executed

# see below for config file structure
COPY shell-conf.yaml /app/
# tell the shell operator where you config file is
ENV SHELL_CONFIG /app/shell-config.yaml
```

You can also set your config in a Kubernetes Config Map and volume it in to reduce how often you have to change your docker image.

Where the `shell-conf.yaml` mentioned above is something like:

```yaml
---
# `boot` is optional
boot:
  # this command is run on boot and is useful for upserting your CRD creation object
  # or any other prep work to be done once when any new pod comes up.
  command: kubectl apply -f /app/mycrd.yaml
  shell: "/bin/bash"
  # Set env vars to be available in the shell
  # This way you can set environment specific items
  # as per a normal 12 factor app
  environment:
    DB_URL: "xxx"

# `watch` is a required key and is an array of watches
watch:
    # The fully qualified API path to a List endpoint
    # This allows you to specify any Kubernetes Objects including native ones such as
    # Pods, Namespaces etc.
  - apiUrl: /apis/my.domain.com/v1/MyCustomObject
    # The command to run - the operator will execute this in a subshell with the default shell
    # it will pipe the Object being updated as a JSON object to stdin.
    command: python myscript.py
    # if this is true the operator will read stderr of the shell and use that to continuosly
    # update the Object's status as it is running.
    # On exit of the subshell if it is a 0 exit code it will expect a json object to have been sent
    # on stdout and will update the CRD's 'status' key with that object.
    # If there is a non 0 exit code it will set the CRD's 'status' as errored and use the last line of
    # stderr as the error message.
    # If this is set to false then the operator will fire and forget - this might be useful
    # for non critical tasks that do not update the Object or for non CRDs that cannot be updated
    # This MUST be false for a non CRD object
    updateObject: true
    # set how many workers the operator should use
    # there may be a need to serialise everything because of race conditions so this can be set to 1
    # or you might want multiple crds to be updated at the same time so
    # want to have 10 or more workers
    # This option is here to constrain the operator to the resources you want to use.
    concurrency: 1
    # The shell to use for the command, this defaults to the default shell
    # but can be overriden
    shell: "/bin/bash"
    # Set env vars to be available in the shell
    # This way you can set environment specific items
    # as per a normal 12 factor app
    environment:
      DB_URL: "xxx"
```

### CRD input

As above, the stdin data that a script will recieve is the following:

```json
{
  "event": "ADDED",
  "object": {

  },
  "ts": 123126423743
}
```

Where:
- event: one of ADDED, MODIFIED or DELETED
- object: is the Object from Kubernetes this event relates to, look at the API server path to see what it might look like for your intended Object
- ts: the unix millisecond timestamp this event occured
