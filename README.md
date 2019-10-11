# Nag

**WORK IN PROGRESS**

[![build status on master](https://api.travis-ci.org/garthk/nag.svg?branch=master)](https://travis-ci.org/garthk/nag)

Nag yourself to fix anything broken, by running Nagios plugins for you:

```sh
go install github.com/garthk/nag
nag list
nag run
```

You'll love it so much, you'll run it every time you open a terminal:

```sh
echo "$(which nag) check" >> ~/.profile
```

Or, you'll just leave it running:

```sh
watch nag check
```

**CHECK _ALL_ THE THINGS!**

![Allie Brosh preparing for her new life as an adult](doc/ALL-THE-THINGS-50.png)

Allie Brosh at her best: [This is Why I'll Never be an Adult][TIWINBA]. If you like it, [buy her book][HAAHB].

[TIWINBA]: http://hyperboleandahalf.blogspot.com.au/2010/06/this-is-why-ill-never-be-adult.html
[HAAHB]: http://hyperboleandahalfbook.blogspot.com.au

## Demo

Try Nag out in Docker:

* `docker run -t -i ubuntu:trusty bash`
* `apt-get update`
* `apt-get install nagios-nrpe-server`
* Ok, you got me. I need to get Travis working with GitHub releases first.

If you prefer `precise`, `apt-get install dialog nagios-nrpe-server` to avoid
spam during installation.

## Audience

I designed Nag for:

* Operators who need to know if they fixed it _now_, not when
  Nagios next decides to poll

* Anyone developing or maintaining Nagios plugins, especially if they are
  [learning Nagios][NFD]

[NFD]:doc/nagios-for-developers.md

## Why the name?

> nag, v: harrass someone to do something

Think of your operators as frustrated Parent Shaped Objects trying to get you
through your morning routine:

* `/usr/local/lib/nagios/plugins/check_brushed_teeth`
* `/usr/local/lib/nagios/plugins/check_cleaned_shoes`
* `/usr/local/lib/nagios/plugins/check_made_bed`
* `/usr/local/lib/nagios/plugins/check_put_away_toys`

Besides, Ubuntu's automatic package finder couldn't find anything:

```
No command 'nag' found, did you mean:
... 8 items, none called nag ...
nag: command not found
```

## Roadmap

* `nag run` executables
* `nag run` commands
* `nag run` in one or all containers via `docker exec`
* `nag run` on another host via `ssh`
* automatic production of NRPE configuration
* automatic production of Nagios configuration

[Plugin API]: http://nagios.sourceforge.net/docs/3_0/pluginapi.html
[exported resources]: https://docs.puppetlabs.com/guides/exported_resources.html
[mines]: http://docs.saltstack.com/en/latest/topics/mine/
[nagiosplugin]: https://pypi.python.org/pypi/nagiosplugin/1.2.2
[NCSA]: https://exchange.nagios.org/directory/Addons/Passive-Checks/NSCA--2D-Nagios-Service-Check-Acceptor/details
[NRPE]: https://exchange.nagios.org/directory/Addons/Monitoring-Agents/NRPE--2D-Nagios-Remote-Plugin-Executor/details
[passive]: https://exchange.nagios.org/directory/Addons/Passive-Checks
