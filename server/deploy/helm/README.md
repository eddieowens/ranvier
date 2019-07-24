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

Then run the Helm install
```bash
helm install --set ssh_key=true ./ranvier
```

`Ranvier` will now be authorized to poll your configuration Git repository! To run Ranvier in non-dev mode, target a 
specific git repo/branch, or to modify where files are stored, see the 
[valid env vars](https://github.com/eddieowens/ranvier/tree/master/server#valid-env-vars).
