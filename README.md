# Ranvier
Ranvier (named after the French neurophysiologist), is a system for managing all of your application's config in a 
centralized, transparent, and fault-tolerant way. It consists of
* A Git repository
* A "Node of Ranvier" server which will poll the specified Git repository
* A Ranvier client which will maintain the configuration within your application

## Setup
1. Create a Git repo with your JSON [config schema](https://github.com/eddieowens/ranvier/wiki/Config-schema-files) 
files.
    * See the [example repo](https://github.com/eddieowens/ranvier-config-example) for a simple example.
1. Setup the [Ranvier server](server/README.md) to point at your created Git repo.
1. Connect to your Ranvier server with one of the clients. Currently supported clients are
    * [Go](client/README.md)
    
## Packages included in Ranvier
Ranvier consists of 5 packages
1. [server](server/README.md)
    * All code for the Ranvier server is stored here. The server handles polling the Git repo which houses your config
    files.
1. [client](client/README.md)
    * The Go client used to communicate with the server. Many more languages soon to be supported.
1. [lang](lang/README.md)
    * The engine for compiling your schema config files to actual config files.
1. [cli](ranvier/README.md)
    * A tool for interacting with Ranvier and supporting local dev while using Ranvier.
1. [commons](commons/README.md)
    * All common code shared amongst all of the packages within Ranvier.