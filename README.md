clementine
==========

Based off of https://github.com/citadel/citadel/tree/master/bastion

### Build

    godep go build
    ./clementine -conf=example/config.json

### Notes

All docker instances (engines) must be accessible over TCP. Additionally, they must all use TLS, and have their certs signed by the same CA. This allows all engines to be managed from a single client cert.

### Examples

An example engine config:

    {
      "id": "boot2docker",
      "addr": "https://192.168.59.103:2376",
      "cpus": 4,
      "memory": 32000,
      "labels": [
        "local",
        "boot2docker"
      ]
    }

An example image config: (note the use of multi as the scheduler, which determines how the container is placed across the cluster)

    {
      "name":"redis:latest",
      "bind_ports": [
        { "proto":"tcp",
          "port":6379,
          "container_port":6379
        }
      ],
      "type": "multi"
    }

#### Add a boot2docker engine

    curl -XPOST '127.0.0.1:8080/apiv1/engine' -d @example/boot2docker_engine.json
    {"ok":true,"data":{"id":"boot2docker","addr":"https://192.168.59.103:2376","cpus":4,"memory":32000,"labels":["local"]}}

#### List engines

    curl '127.0.0.1:8080/apiv1/engines'
    {"ok":true,"data":[{"id":"boot2docker","addr":"https://192.168.59.103:2376","cpus":4,"memory":32000,"labels":["local"]}]}

#### List Containers on Engine

    curl '127.0.0.1:8080/apiv1/engine/boot2docker/containers'
    {"ok":true,"data":[{"id":"f810083c6bfe204997cc17352cdf860fee757fcc03a2c29100ef81ebc09dc054","name":"/sad_hodgkin","image":{"name":"redis:latest","entrypoint":["/entrypoint.sh"],"environment":{"REDIS_DOWNLOAD_SHA1":"913479f9d2a283bfaadd1444e17e7bab560e5d1e","REDIS_DOWNLOAD_URL":"http://download.redis.io/releases/redis-2.8.17.tar.gz","REDIS_VERSION":"2.8.17"},"hostname":"f810083c6bfe","type":"multi","labels":[""],"bind_ports":[{"proto":"tcp","host_ip":"0.0.0.0","port":6379,"container_port":6379}],"volumes":["/data"],"restart_policy":{},"network_mode":"bridge"},"engine":{"id":"boot2docker","addr":"https://192.168.59.103:2376","cpus":4,"memory":32000,"labels":["local"]},"state":"running","ports":[{"proto":"tcp","host_ip":"0.0.0.0","port":6379,"container_port":6379}]}]}

#### Destroy a running container

    curl -v  -XDELETE '127.0.0.1:8080/apiv1/engine/boot2docker/container/f810083c6bfe204997cc17352cdf860fee757fcc03a2c29100ef81ebc09dc054'
    {"ok":true,"data":"Container container /sad_hodgkin Image redis:latest Engine boot2docker destroyed on engine engine boot2docker addr https://192.168.59.103:2376"}

#### Run a container on the cluster

    curl '127.0.0.1:8080/apiv1/container' -d @example/redis_image.json
    {"ok":true,"data":{"id":"c16b1738fad321310e1bafe39fb09f2800dcdedadaab708362f8b180add98f8b","image":{"name":"redis:latest","type":"multi","bind_ports":[{"proto":"tcp","port":6379,"container_port":6379}],"restart_policy":{}},"engine":{"id":"boot2docker","addr":"https://192.168.59.103:2376","cpus":4,"memory":32000,"labels":["local"]},"ports":[{"proto":"tcp","host_ip":"0.0.0.0","port":6379,"container_port":6379}]}}

## TODO

* Allow listing of ALL containers (and handle failure of engines gracefully)
* Rethink redis schema
* Provide API to update status of a engine
* Store images in a manner that they can be recalled easily (i.e. schedule a saved redis image config)
* Allow scaling of an instance (everything under http://shipyard-project.com/docs/api/ honestly)
* AUTHENTICATION
