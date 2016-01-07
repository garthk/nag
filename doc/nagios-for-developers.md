# Nagios for Developers

## TL;DR

* Write a short script that checks something and exits appropriately
* Check your script with [`nag`](../)
* Drop a `.cfg` file into `/etc/nagios/nrpe.d/` naming and calling your script
* Check your command with `nag`
* Let your operators know the good news
* Repeat until you **CHECK ALL THE THINGS!**

![Allie Brosh preparing for her new life as an adult](ALL-THE-THINGS-50.png)

Allie Brosh at her best: [This is Why I'll Never be an Adult][TIWINBA]. If you like it, [buy her book][HAAHB].

[TIWINBA]: http://hyperboleandahalf.blogspot.com.au/2010/06/this-is-why-ill-never-be-adult.html
[HAAHB]: http://hyperboleandahalfbook.blogspot.com.au)

## Is Nagios Too Hard?

Any of this sound familiar?

You decide to start with Nagios Core. The documentation seems to never end.
You quickly discover the configuration is brittle. If you get anything wrong,
Nagios stops working. Manual configuration work gives you plenty of
opportunities to get it wrong. Chaos ensues.

You buy Nagios XI for the web interface. It's easiest to configure Nagios to
check your server from the outside, so that's what you do. If it's down, you
know it's down, but not *why* it's down. You open your [runbook] to the
troubleshooting steps and start pasting commands into your terminal. The
runbook is never comprehensive enough or up to date. Chaos continues ensuing.

[runbook]: https://en.wikipedia.org/wiki/Runbook

You automate some of the troubleshooting steps and deploy them along with your
software. Emboldened, you try Puppet [exported resources] or Salt [mines] to
gather monitoring needs and update your Nagios configuration to run those
tools. It works — BWA HA HA! — but your Puppet and Salt configurations become
hard to work with. Your developers and operators, already overwhelmed, refuse
to help. They want to stop automating entirely and go back to the runbooks,
because it's easier. Chaos continues ensuing.

Before you moan [`#monitoringsucks`][monitoringsucks] on Twitter and give up:

[monitoringsucks]: http://www.kartar.net/2013/01/monitoring-sucks---a-rant/

It doesn't have to be that hard. You can take one step at a time, and get
benefit from each step. It's this easy: **start at the other end**. Start on
your server. Not the Nagios server, *your* server. The one doing the work.

Eat dessert first.

## Nagios, Backwards

* Configure Nagios NRPE commands on your servers
* Run them with [Nag](../) from the terminal during maintenance and troubleshooting
* Install NRPE
* Montior the commands from Nagios

### Basic Nagios Plugin

The first step is to write a Nagios [plugin][Plugin API]. Nagios plugins are
simple scripts or executables that check something. They let Nagios know their
status through their [exit status] and some short text output.

The simplest Nagios plugin rites a short line of human-readable text and exits with an [exit status] corresponding to a service status:

* `0`: `OK`
* `1`: `WARNING`
* `2`: `CRITICAL`
* `3`: `UNKNOWN`

Let's say you need to make sure that `service nginx status` is happy. Paste
this into `/usr/local/bin/check_nginx_service` and mark it executable:

```sh
#!/bin/sh
service nginx status || exit 2
```

The `check_nginx_service` script above will exit `OK` if the `service` command
exits with status 0. Otherwise, it will exit `CRITICAL`. Its output will be
whatever `service` outputs.

No, really, it's that easy.

If you feel bold, now, upgrade to a more advanced version:

```sh
#!/bin/sh
service nginx status 1>&2  # pipe stdout to stderr so Nagios doesn't see it
case $? in
  0) echo OK;    exit 0;;  # Hooray!
  1) echo DEAD;  exit 2;;  # PID file exist, but process dead
  3) echo DEAD;  exit 2;;  # not running or PID file empty
  *) echo WTF;   exit 3;;  # no idea what's going on
esac
```

The new `check_nginx_service` script gives its own text output. It also exits
`UNKNOWN` if the `service` output isn't helpful.

Run your plugin via `nag` to check its behaviour:

* `nag run /usr/local/bin/check_nginx_service`

### Basic NRPE Command

Paste this into `/etc/nagios/nrpe.d/site.cfg`:

```plain
command[check_nginx_service]=/usr/local/bin/check_nginx_service
```

Then:

* Check your command's configuration with `nag run check_nginx_service`
* Tell the operators the good news
* Enjoy checking ALL THE THINGS when you `nag run` without arguments

### Using Default Plugins

If you install NRPE's default plugins, you can call them with arguments from
your NRPE command configuration without writing your own plugin:

* Ubuntu-style: `apt-get install nagios-nrpe-server`
* RedHat-style: `yum install nrpe nagios-plugins-all --enablerepo=epel`

You'll end up with handy plugins in `/usr/lib/nagios/plugins` or
`/usr/lib64/nagios/plugins`.

Need to check your web server is responding on `localhost:8001` but can't
reach it from your Nagios server? Paste this into
`/etc/nagios/nrpe.d/site.cfg`:

```plain
command[check_internal_server]=/usr/lib/nagios/plugins/check_http http://localhost:8001
```

### Using existing test commands

Use `nag run -S` to turn any executable test into a Nagios plugin.

```plain
command[check_log_directory]=/usr/local/bin/nag run -S -- test -d /var/log/apache2
```

In sensitive mode, `nag` will:

* Run the command
* Treat its exit status as `CRITICAL` if non-zero
* Truncate its output at Nagios' default `MAX_PLUGIN_OUTPUT_LENGTH` (4KB)
* If output is empty, replace it with `OK` or `CRITICAL`
* Ensure the output ends with a newline

`man test` for a quick refresher on all your new tests.

### Arguments

TODO: explain why the ideal number of arguments is zero, and how to re-use
your application's configuration to check its dependencies.

### Docker

TODO: explain you can run `nag` inside your Docker containers with `docker
exec`, keeping them small by including only `nag`, a minimal or no
`nrpe.cfg`, and a handful of plugin scripts.

```plain
command[check_container]=docker exec containername nag --all
```

### Deployment

TODO: Persuade devops teams to include NRPE command files when they automate
deployment.

### Packaging

TODO: Persuade packagers to include NRPE command files when they package
software.

### Minimizing Nagios Configuration

TODO: Explain how to run `nag` via NRPE to run all the other commands and
summarise the status.

### Skipping Nagios Entirely

TODO: Explain how to hook `nag` up to other popular monitoring software.

[exit status]: http://en.wikipedia.org/wiki/Exit_status
[Plugin API]: https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/pluginapi.html
[exported resources]: https://docs.puppetlabs.com/guides/exported_resources.html
[mines]: http://docs.saltstack.com/en/latest/topics/mine/
[NCSA]: https://exchange.nagios.org/directory/Addons/Passive-Checks/NSCA--2D-Nagios-Service-Check-Acceptor/details
[NRPE]: https://exchange.nagios.org/directory/Addons/Monitoring-Agents/NRPE--2D-Nagios-Remote-Plugin-Executor/details
[passive]: https://exchange.nagios.org/directory/Addons/Passive-Checks
