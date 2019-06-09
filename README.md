# Mannequin

![logo](./internal/misc/mannequin.png "Mannequin")

_A mannequin (also called a manikin, dummy, lay figure or dress form) is an often articulated doll used by artists, tailors, dressmakers, windowdressers and others especially to display or fit clothing. The term is also used for life-sized dolls with simulated airways used in the teaching of first aid, CPR, and advanced airway management skills such as tracheal intubation and for human figures used in computer simulation to model the behavior of the human body. During the 1950s, mannequins were used in nuclear tests to help show the effects of nuclear weapons on humans._

## Why do I need this?
Basically this is a thing to simplify your local development of microservices that suppose to run in k8s.
Especially it should become handy if you have a lot of dependencies or you need to work in multiple services at the same time.
Idea is that you can just write some configs for your local projects, configure it's dependencies (both services like mysql, beanstalk or redis; or your local codebases/other projects that are registered in this thing( and then it would autodetect if you have minikube or docker-for-desktop, build your image, deploy it there using helm. You can enable `watch` that would build new image and redeploy it locally on codechange.

## Features
1. Works with Kubernetes
- Before deployment checks if environment matches the defined one in configuration
2. Configurable from file
- [X] configuration per project set
- [X] configuration per repo
3. ~Compatible with Gitlab CI pipeline~
- ~able to understand steps from CI pipeline~
4. Supports cloud emulators
- pubsub
5. [X] Has own client called mnqnctl
- [X] Used to initialize new project set
- Used to start listener with autodeploy for a project
6. [X] Compatible with helm for deployments
7. External debuggers support
- Go Delve
8. Supports local k8s
- Minikube
- Docker for desktop
9. Support for service dependencies
- Codebased dependencies (other local projects)
- Service based dependencies (mysql, postgres, redis, etc)
10. Local ingress support

### TODO:
- Write HOWTO for minikube and DFD
- Write troubleshooting
- Add support for preselecting context based on local config

## mnqnctl

### Init
```
mnqnctl init
```

Starts a questionaire to set up your local project environment.
Registers project to the list of deployable codebases.

### Implode

```
mnqnctl implode
```

Implodes previously made global configuration.

### Deploy

```
mnqnctl deploy
```

Deploys project with the configuration in the same folder via selected kubernetes context.
Lints helm charts by default.

```
mnqnctl deploy latest
```

Deploys the project along with the latest code version of the dependencies with type project.

### Watch

```
mnqnctl watch
```

Deploys project with the configuration in the same folder on code changes via selected kubernetes context.
