# Ranvier CLI
The CLI tool for interacting with Ranvier. This CLI houses the exact same compiler code that the 
[server](../server/README.md) uses giving you the ability to compile, lint, and validate your config schema files 
locally or through your automated build system. 

It also facilitates local development by allowing developers to run 
their application code locally without the need for starting the Ranvier server (although that's very 
[simple](../server/deploy/docker/README.md)). 

## Installation
```bash
go get -u github.com/eddieowens/ranvier/ranvier
```

## Usage
```bash
ranvier compile [config schema filepath]
```
This command will compile the specific config schema file to a singular config file for your application to consume

### Example
Let's say you want to develop a feature on the `users` service which uses Ranvier and you need its config file. First
clone your Git repo storing all of the config schema files (let's say it's one directory above your current) and simply 
run
```bash
ranvier compile -r ../my-config-repo users.json
```
A `users.json` file should appear in your current directory with all of the needed config.

### Getting help
All commands on the cli are equipped with extensive help docs. Simply run
```bash
ranvier help compile
```
To see the help docs as well as available flags.
