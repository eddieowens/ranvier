## Node of Ranvier
This directory houses all of the code to run the Ranvier server.

### Running

#### Natively
```bash
go run main.go
```

#### Docker
Ranvier needs either an SSH key or a username/password to poll your git repository. If none
are provided, the server will fail to start

Run with SSH
```bash
docker run -p 8080:8080 -e RANVIER_GIT_SSHKEY=/.ssh/id_rsa -v ~/.ssh/id_rsa:/.ssh/id_rsa edwardrowens/ranvier
```

Run with username/password
```bash
docker run -p 8080:8080 -e RANVIER_GIT_USERNAME=eddieowens -e RANVIER_GIT_PASSWORD=<password> edwardrowens/ranvier
```

### Configuration
The configuration for the server is handled via files in `config`. The `config.yml` file
is loaded first and the subsequent `ENV` specific file is merged in over taking precedence
in the case of a key collision.

#### Configure via environment variables
All of the env vars are overrideable via environment variables and their name is automatically
allocated. Every env var is prefixed with `RANVIER_` and every underscore represents a
traversal down the object tree.

#### Valid Env Vars
| Key                        | Description                                                                                               | Default                                                |
|----------------------------|-----------------------------------------------------------------------------------------------------------|--------------------------------------------------------|
| `ENV`                      | The environment that the server will operate in. Affects which of the `config/<ENV>.yml` files is loaded. | `dev`                                                  |
| `SERVER_PORT`              | The port that the server will run on.                                                                     | `8080`                                                 |
| `GIT_REMOTE`               | The remote git URL that contains your configuration.                                                     | `git@github.com:eddieowens/ranvier-config-example.git` |
| `GIT_BRANCH`               | The branch that the server will watch for changes.                                                        | `master`                                               |
| `GIT_DIRECTORY`            | The directory that the repository will be cloned into.                                                    | `./ranvier-config-example`                             |
| `GIT_POLLINGINTERVAL`      | The interval (in seconds) to poll the git remote for changes                                              | `10`                                                   |
| `GIT_USERNAME`             | The username used to auth the git remote. Ignored if `GIT_PASSWORD` is not also set.                      |                                                        |
| `GIT_PASSWORD`             | The password used to auth the git remote. Ignored if `GIT_USERNAME` is not also set                       |                                                        |
| `GIT_SSHKEY`               | The full filepath to the private SSH key used to auth the git remote.                                     |                                                        |
| `COMPILER_OUTPUTDIRECTORY` | The directory that the compiled configuration files will live on the server.                              | `./output`                                             |