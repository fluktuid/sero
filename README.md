# Sero

Scale to zero in your cluster - but simple.

`Sero` works as a small interceptor watching the connections.
So whenever a connection to `sero` is opened it will try to forward your connection.
If your Application is scaled to zero, `sero` will scale it to one.
The rest will then be done by the awesome Kubernetes Autoscaler.


## How to use

1. Set up Sero
1. Install sero for the desired component
1. Send your Traffic to sero
1. 🎉


### Example

Imagine there is an app 'a' that is reached in the cluster under 'a-svc:80'.
The deployment of the app is 'a-deploy'.

In this case, a possible configuration looks like this:

``` yaml
target:
  host: 'a-svc'
  port: 80
  protocol: tcp
  deployment: a-deploy
  timeout:
    forward: 200 # maximum waiting time when forwarding a request
    scaleUP: 3000 # maximum waiting time after an 'scale up' event
```

## To concider

- Sero has not yet been tested for production use.
  _(If you did, feel free to contact the project team.)_
- Ideally, sero and target pod should run on the same node.
  Affinities can help here.
- Sero needs some resources itself to trade connections.
  The requests / limits or scalers must therefore be tailored to the load growth of your application.
