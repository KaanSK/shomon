<p align="center">
  <img src="images/logo.png" />
</p>

<p align="center">
  <img src="https://goreportcard.com/badge/github.com/KaanSK/shomon" />
  <img src="https://img.shields.io/github/downloads/kaansk/shomon/total?color=%233ABE25&label=Release%20Downloads" />
  <img src="https://img.shields.io/docker/pulls/kaansk/shomon?color=%233ABE25&label=DockerHub%20Pulls" />
</p>


<p align="center">
ShoMon is a Shodan alert feeder for TheHive written in GoLang. With version 2.0, it is more powerful than ever!
</p>


# Functionalities
* Can be used as Webhook OR Stream listener
    * Webhook listener opens a restful API endpoint for Shodan to send alerts. This means you need to make this endpoint available to public net
    * Stream listener connects to Shodan and fetches/parses the alert stream
* Utilizes [shadowscatcher/shodan](https://github.com/shadowscatcher/shodan) (fantastic work) for Shodan interaction.
* Console logs are in JSON format and can be ingested by any other further log management tools
* CI/CD via Github Actions ensures that a proper Release with changelogs, artifacts, images on ghcr and dockerhub will be provided
* Provides a working [docker-compose file](docker-compose.yml) file for TheHive, dependencies
* Super fast and Super mini in size
* Complete code refactoring in v2.0 resulted in more modular, maintainable code
* Via conf file or environment variables alert specifics including tags, type, alert-template can be dynamically adjusted. See [config file](conf.yaml).
* Full banner can be included in Alert with direct link to Shodan Finding.

    ![Alert example](images/alert.png)
* IP is added to observables

    ![Observable example](images/observable.png)

# Usage
* Parameters should be provided via ```conf.yaml``` or environment variables. Please see [config file](conf.yaml) and [docker-compose file](docker-compose.yml)
* After conf or environment variables are set simply issue command: 

    `./shomon`

## Notes
* Alert reference is first 6 chars of md5("ip:port")
* Only 1 mod can be active at a time. Webhook and Stream listener can not be activated together.

# Setup & Compile Instructions
## Get latest compiled binary from releases
1. Check [Releases](https://github.com/KaanSK/shomon/releases/latest)  section.

## Compile from source code
1. Make sure that you have a working Golang workspace.
2. `go build .`
    * `go build -ldflags="-s -w" .` could be used to customize compilation and produce smaller binary.

## Using Public Container Registries
1. Thanks to new CI/CD integration, latest versions of built images are pushed to ghcr, DockerHub and can be utilized via:
    * `docker pull ghcr.io/kaansk/shomon`
    * `docker pull kaansk/shomon`

## Using [Dockerfile](Dockerfile)
1. Edit [config file](conf.yaml) or provide environment variables to commands bellow
2. `docker build -t shomon .`
3. `docker run -it shomon`

## Using [docker-compose file](docker-compose.yml)
1. Edit environment variables and configurations in [docker-compose file](docker-compose.yml)
2. `docker-compose run -d`

# Credits
* Logo Made via LogoMakr.com
* [shadowscatcher/shodan](https://github.com/shadowscatcher/shodan) 
* [Dockerfile Reference](https://www.cloudreach.com/en/resources/blog/cts-build-golang-dockerfiles/) 
* Release management with [GoReleaser](https://goreleaser.com)
