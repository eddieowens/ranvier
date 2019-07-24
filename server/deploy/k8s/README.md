# K8s Definitions for Ranvier
Contains all of the files required to run a node of `Ranvier` in dev mode in k8s. Requires a secret is created 
containing either your git password or an SSH key for git.

## Git access with username/password
The `Ranvier` `Deployment` found in `deploy.yml` by default assumes a `Secret` is present in the same namespace 
containing the password as well as an associated username.

1\. To create a secret for your password,
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: git-access
type: Opaque
data:
  password: <base64 encoded password here>
```

2\. Then add a reference to that secret in `deploy.yml`
```yaml
...
  env:
    - name: RANVIER_GIT_USERNAME
      value: <your username>
    - name: RANVIER_GIT_PASSWORD
      valueFrom:
        secretKeyRef:
          key: password
          name: git-access
...
```

`Ranvier` will now be authorized to poll your configuration Git repository!

## Git access with SSH key
To use an SSH key, 
1. In `deploy.yml` remove the environment variables referencing `RANVIER_GIT_USERNAME` and `RANVIER_GIT_PASSWORD` and 
replace them with an `env` key of `RANVIER_GIT_SSHKEY` and a value of `/.ssh/id_rsa`.
```yaml
...
  env:
    - name: RANVIER_GIT_SSHKEY
      value: /.ssh/id_rsa
..
```
2\. Create a `volume` in the `Deployment`
```yaml
...
  volumes:
    - name: ssh-key
      secret:
        secretName: git-access
        items:
          - key: ssh-key
            path: id_rsa
...
```
3\. Create a `volumeMount` for the container
```yaml
...
  volumeMounts:
    - mountPath: /.ssh/id_rsa
      subPath: id_rsa
      name: ssh-key
      readOnly: true
...
```

4\. Create the `Secret` for your node of `Ranvier`

By yaml definition
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: git-access
type: Opaque
data:
  ssh-key: <base64 encoded contents of your SSH key here>
```
Or by `kubectl`
```bash
kubectl create secret generic git-access --from-file=ssh-key=$HOME/.ssh/id_rsa
```

`Ranvier` will now be authorized to poll your configuration Git repository! To run Ranvier in non-dev mode, target a 
specific git repo/branch, or to modify where files are stored, see the 
[valid env vars](https://github.com/eddieowens/ranvier/tree/master/server#valid-env-vars).
