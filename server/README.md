## Node of Ranvier
This directory houses all of the code to run the Ranvier server.

### Running

#### Natively
```bash
go run main.go
```

#### Containerized
* [Docker](deploy/docker/README.md)
* [Helm](deploy/helm/README.md)
* [Kubernetes](deploy/k8s/README.md)

### Configuration
The configuration for the server is handled via files in `config`. The `config.yml` file
is loaded first and the subsequent `ENV` specific file is merged in taking precedence
in the case of a key collision.

#### Configure via environment variables
All of the config values are overrideable via environment variables and their name is automatically
allocated. Every env var is prefixed with `RANVIER_` and every underscore represents a
traversal down the object tree.

#### Valid Env Vars
| Key                        | Description                                                                                                               | Default                                                |
|----------------------------|---------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------|
| `ENV`                      | The environment that the server will operate in. Affects which of the `config/<ENV>.yml` files is loaded.                 | `dev`                                                  |
| `SERVER_PORT`              | The port that the server will run on.                                                                                     | `8080`                                                 |
| `GIT_REMOTE`               | The remote git URL that contains your configuration.                                                                      | `git@github.com:eddieowens/ranvier-config-example.git` |
| `GIT_BRANCH`               | The branch that the server will watch for changes.                                                                        | `master`                                               |
| `GIT_DIRECTORY`            | The directory that the repository will be cloned into.                                                                    | `./ranvier-config-example`                             |
| `GIT_POLLINGINTERVAL`      | The interval (in seconds) to poll the git remote for changes                                                              | `10`                                                   |
| `GIT_USERNAME`             | The username used to auth the git remote. Ignored if `GIT_PASSWORD` is not also set.                                      |                                                        |
| `GIT_PASSWORD`             | The password used to auth the git remote. Ignored if `GIT_USERNAME` is not also set                                       |                                                        |
| `GIT_SSHKEY`               | The full filepath to the private SSH key used to auth the git remote.                                                     |                                                        |
| `COMPILER_OUTPUTDIRECTORY` | The directory that the compiled configuration files will live on the server.                                              | `./output`                                             |
| `LOG_LEVEL`                | The case-insensitive logging level of the server. Valid levels are trace, debug, info, warn, warning, error, fatal, panic | `info`                                                 |
| `LOG_TIMEFORMAT`           | The [golang](https://gobyexample.com/time-formatting-parsing) time format string.                                         | `RFC3339`                                              |