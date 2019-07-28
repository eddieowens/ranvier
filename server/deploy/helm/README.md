# Helm chart for Ranvier
Contains all of the files required to run a node of `Ranvier` in dev mode in k8s via Helm. Requires a secret is created 
containing either your git password or an SSH key for git. 

## Using username/password
To install `Ranvier` via password auth, run
```bash
helm install --set password=<your git password> --set username=<your git username> ./ranvier
```
## Using SSH
Due to limitations in Helm, it is not possible to read an SSH key from a local directory
into a k8s `Secret`. To run `Ranvier` using SSH git auth, create a `Secret` like the following
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: git-access
type: Opaque
data:
  ssh-key: <base64 encoded contents of your private SSH key for git>
```
Or create it via `kubectl`
```bash
kubectl create secret generic git-access --from-file=ssh-key=$HOME/.ssh/id_rsa
```
If you don't have an `id_rsa` key in your `~/.ssh` directory, see 
[this](https://help.github.com/en/articles/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) on how to make one.

Then run the Helm install
```bash
helm install --set ssh_key=true ./ranvier
```

`Ranvier` will now be authorized to poll your config's Git repository! By default, Ranvier will sync with the 
[example repo](https://github.com/eddieowens/ranvier-config-example). To run Ranvier in non-dev mode, target **your** 
git repo/branch, or to modify where files are stored, see the
[valid env vars](https://github.com/eddieowens/ranvier/tree/master/server#valid-env-vars).

## Available Values
| Key                      | Type     | Description                                                                                                                                                                                                                  | Default                       |
|--------------------------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------|
| `ssh_key`                | `bool`   | If set to `true`, the server will run with the assumption that there is a `secret` with the name of the `secret.git_access.name` value and a field name of `ssh-key`. The secret is stored in the pod's `/.ssh/id_rsa` file. |                               |
| `username`               | `string` | The Git username for the target repo. Required if the `password` value is set.                                                                                                                                               |                               |
| `password`               | `string` | The password for the target repo. Required if the `username` value is set.                                                                                                                                                   |                               |
| `secret.git_access.name` | `string` | The name of the secret which stores the contents of the SSH key used to auth into the Git repo.                                                                                                                              | `git-access`                  |
| `container.port`         | `int`    | The port number that the server will listen on.                                                                                                                                                                              | `8080`                        |
| `container.image.name`   | `string` | The name of the Docker image to use for the Ranvier server.                                                                                                                                                                  | `edwardrowens/ranvier-server` |
| `container.image.tag`    | `string` | The tag used for the Docker image of the Ranvier server.                                                                                                                                                                     | `latest`                      |
