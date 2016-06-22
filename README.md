## RabbitHook Builder

[RabbitHook][rh] provides endpoints to handle webhook events from services like [Dockerhub][dw] or [Github][gw], and then publishes those events over AMQP.

Rabbit Builder is an agent for automated `docker` builds using [RabbitHook][rh] events.

`rhbuilder` is a long running AMPQ client subscribed to events on the `rabbithook.dockerhub` topic. Any time we push an image to `Dockerhub`, an event is published and `rhbuilder` will then try to execute a build job for each event.

---

### Adding jobs

To create a new job, you just need to place an executable script following some simple conventions:

```
base-path + / + dockerhub.owner + / + dockerhub.repository_name
```

* `base-path`: By default is set to `/usr/local/opt/rhbuilder`
* `dockerhub.owner`: Usually matches your Dockerhub username
* `dockerhub.repository_name`: Matches your Dockerhub project's name

If you have a project on Dockerhub, say [goliatoe/hello-rabbit][hr]. You need to create a directory named `goliatone` under `/usr/local/opt/rhbuilder` and place an executable script named `hello-rabbit`.

Note that in order to work, [goliatoe/hello-rabbit][hr] has a webhook pointing to a running [RabbitHook][rh] instance.

- topic:
    - `rabbithook.dockerhub`


Sample Event Payload:

```json
{
    "push_data": {
        "pushed_at": 1466026222,
        "images": [],
        "tag": "latest",
        "pusher": "goliatone"
    },
    "callback_url": "https://registry.hub.docker.com/u/goliatone/hello-rabbit/hook/4hac5445ga3f345ajeje31abc34ie06rip/",
    "repository": {
        "status": "Active",
        "description": "Hello rabbit!",
        "is_trusted": false,
        "full_description": "",
        "repo_url": "https://hub.docker.com/r/goliatone/hello-rabbit",
        "owner": "goliatone",
        "is_official": false,
        "is_private": false,
        "name": "hello-rabbit",
        "namespace": "goliatone",
        "star_count": 0,
        "comment_count": 0,
        "date_created": 1461793462,
        "repo_name": "goliatone/hello-rabbit"
    }
}
```

## Daemon

This repository also provides a set of scripts to run `rhbuilder` using `init.d` or `upstart`. If you want to run this on a Raspberry Pi use the `init.d` scripts, else if possible use `upstart`.

For more information on how to use the boot scripts, refer to the readme file for the [init.d][initd] or [upstart][upstart]

[initd]:https://github.com/goliatone/rabbithook-builder/blob/master/conf/init.d/README.md
[upstart]:https://github.com/goliatone/rabbithook-builder/blob/master/conf/ubuntu/README.md

### Install

You can use the installation scripts to install the daemon on your machine.

```
curl -sL https://raw.githubusercontent.com/goliatone/rabbithook-builder/master/conf/init.d/install | bash
```

### Uninstall
```
curl -sL https://raw.githubusercontent.com/goliatone/rabbithook-builder/master/conf/init.d/uninstall | bash
```


[rh]:https://github.com/goliatone/rabbithook
[gw]:https://developer.github.com/webhooks/
[dw]:https://docs.docker.com/docker-hub/webhooks/
[hr]:https://hub.docker.com/r/goliatone/hello-rabbit/
