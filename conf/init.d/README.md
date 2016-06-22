## init.d

Set of scripts and files needed to run `rhbuilder` as a deamon using `init.d`.

### Files

#### rhbuilder
This file should be placed in `/etc/init.d/rhbuilder`.  This service is used start / stop rhbuilder at boot. It will also load configuration options

#### rhbuilder.conf
Logwatch [configuration][lr] file for logrotate. Rename the file to `rhbuilder`, and place it inside `/etc/logrotate.d/`.

You can modify this file to suit your needs.

For instance, if you want to make sure your log file is created, replace `nocreate` with:
```
create 640 pi pi
```

#### rhbuilder.tpl
Rename the file to `rhbuilder`, and place it inside `/etc/default/rhbuilder`.

Edit the variables to suit your needs:

```
# Required AMQP URL
RHBUILDER_AMQP_URL=amqp://guest:guest@localhost:5672/

# AMQP Exchange
RHBUILDER_AMQP_EXCHANGE=rabbithook

# AMPQ Topic, by default one of rabbithook.[dockerhub|github|travis]
RHBUILDER_AMQP_TOPIC=rabbithook.dockerhub

# Path where "rhbuilder" looks for job files
RHBUILDER_JOBS_PATH=/usr/local/opt/rhbuilder
```

[lw]:http://debianhelp.co.uk/logwatch.htm
