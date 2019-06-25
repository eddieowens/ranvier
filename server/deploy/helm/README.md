# Helm chart for Ranvier
`Ranvier` requires Git access authorized either by SSH or by username/password. 

## Using password
To install `Ranvier` via password auth, run
```bash
helm install --set password=<your git password> ./ranvier
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
_Note: Both the `name` and `ssh_key` fields are important and should not be changed_

Then run the Helm install
```bash
helm install --set ssh_key=true ./ranvier
```
