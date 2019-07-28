# Running Ranvier in Docker
Ranvier needs either an SSH key or a username/password to poll your git repository. If neither
are provided, the server will fail to start

## SSH
```bash
docker run -p 8080:8080 -e RANVIER_GIT_SSHKEY=/.ssh/id_rsa -v ~/.ssh/id_rsa:/.ssh/id_rsa edwardrowens/ranvier-server
```
This will utilize your `~/.ssh/id_rsa` SSH key to poll the Git repo. If you do not have one, generate one 
[like so](https://help.github.com/en/articles/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent).

## Username/password
```bash
docker run -p 8080:8080 -e RANVIER_GIT_USERNAME=eddieowens -e RANVIER_GIT_PASSWORD=<password> edwardrowens/ranvier
```
This method will use your Git username and password to allow Ranvier to poll your Git repo.

Your server is now available on port 8080 and by default will poll the 
[example repo](https://github.com/eddieowens/ranvier-config-example). Try sending it a 
[query](http://localhost:8080/api/config/users-staging).

To not use the example repo, modify the port, etc, see the 
[valid env vars](https://github.com/eddieowens/ranvier/tree/master/server#valid-env-vars) to configure the server. 