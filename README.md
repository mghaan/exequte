[![Chat on Matrix](https://matrix.to/img/matrix-badge.svg)](https://matrix.to/#/#exequte:matrix.org)

# exeQute

A simple modular cross-platform background service to remotely control your device via MQTT writen in Go.

> This is still a work in progress. However, it is possible to compile for various platforms and use regularly.

### Contact

If you have any comments or ideas, you can drop a chat message in [#exequte:matrix.org](https://matrix.to/#/#exequte:matrix.org) via [Matrix](https://www.matrix.org).

## Requirements

* MQTT server

This program runs on all platforms supported by Go. For conviecence precompiled binary packages are provided for Linux, Windows and OS X.

## Building

You can use **make** to speed up the build process:

    make clean
    make all

Compile for specific platform (check Makefile for all targets):

    make linux-amd64

## Installing

Copy **exequte** binary, plugins and configuration file to your device or install exeQute via provided distribution packages.

## Running

Modify the **exequte.json** configuration file to your needs, then start exeQute and point it to your config:

    exequte --config <path>

Distribution packages contain systemd unit you can use to control exeQute:

    systemctl start exequte

## Configuration

All configuration is done in the **exequte.json** config file. If you installed exeQute via distribution package, the file location is **/etc/exequte.json**.

    {
        "system": {
            /**
             * This will enable debug mode. Useful to see more details when an error occurs.
             * /
            "debug": false,
            /**
             * ExeQute will write any log events into this file if set. File log is disabled if empty.
             * The log is always cleared when exeQute is started.
             * /
            "log": ""
        },
        /**
        * Set connection details to your MQTT server here.
        * /
        "mqtt": {
            /**
            * MQTT server host/IP address.
            * /
            "host": "127.0.0.1",
            /**
            * MQTT server port.
            * /
            "port": 1883,
            /**
            * Enable or disable SSL while connecting to MQTT server.
            * /
            "ssl": false,
            /**
            * MQTT server username.
            * /
            "username": "mqtt",
            /**
            * MQTT server password.
            * /
            "password": "secret",
            /**
            * Client ID to present to MQTT. This is useful if you have multiple exeQute instances.
            * Change to an unique ID for each of your instance. It is the prefix to use when subscribing topics.
            * /
            "client": "exequte"
        },
        /**
        * Load plugins and libraries. Go supports external plugins, this allows to extend exeQute with additional actions.
        * Unfortunately, Go supports plugins in Linux only for now. Therefore there are transitional "libraries"
        * to mimic this feature - it's just a separate code so users can customize exeQute behavior.
        "plugins": [
            {
                /**
                * Full path to external plugin. This plugin does nothing - it serves just as a placeholder
                * to show the capability to load plugins on runtime.
                * /
                "plugin": "/usr/local/exequte/plugins/dummy.so"
            },
            {
                /**
                * This is an internal library to execute arbitrary commands on the device via MQTT. You can define
                * a set of actions to execute and its parameters.
                * ExeQute will subscribe topic "exequte/system/run" to listen for actions/commands to execute.
                * /
                "plugin": "run",
                "config": [
                    {
                        /**
                        * Name of the action. Send "ping" in the topic "exequte/system/run" to execute.
                        * /
                        "alias": "ping",
                        /**
                        * Execute the command and parameters exactly as specified here. When we receive "ping"
                        * in the topic "exequte/system/run" then execute ping localhost (it has no meaning to ping
                        * localhost, this just shows how to run arbitrary commands).
                        * /
                        "script": "/usr/bin/ping 127.0.0.1",
                        /**
                        * Should we accept parameters? When set to false, we will use the first parameter as alias.
                        * Anything else is ignored. If set to true, you can specify command and parameters. The
                        * first parameter is the alias, the rest of parameters will be passed to the command to run.
                        "params": false
                    },
                    {
                        /**
                        * Define another alias.
                        * /
                        "alias": "pingmore",
                        /**
                        * Define what command to run when we receive "pingmore" in "exequte/system/run" topic.
                        * /
                        "script": "/usr/bin/ping",
                        /**
                        * We will accept parameters. This means when we receive "pingmore 192.168.1.1", we will
                        * execute "/usr/bin/ping 192.168.1.1".
                        * Remember exeQute is not doing any sanity or validation checks, it just passes the
                        * parameters to the command (script) you specify!
                        "params": true
                    }
                ]
            },
            {
                /**
                * Check process by its name and return zero when running or return a non-zero value when not running.
                * This action is platform specific and uses tools provided by the operating system itself.
                * The topic used to return values is "exequte/system/alive/[process-to-check]".
                * /
                "plugin": "alive",
                "config": [
                    {
                        /**
                        * How often in minutes to check if process is running.
                        * /
                        "interval": 10,
                        /**
                        * Topic where we should publish our findings.
                        * Use "exequte/system/alive/terminal" in this example.
                        * /
                        "topic": "terminal",
                        /**
                        * Process name to check.
                        * /
                        "process": "xterm"
                    },
                    /**
                     * Of course you can continue and specify more actions...
                     * /
                    {
                        "interval": 60
                        "topic": "hourly"
                        "process": "meaningoflife"
                    },
                    {
                        ...
                    }
                ]
            },
            {
                /**
                * Execute command and return its output.
                * The topic used to return values is "exequte/system/status/[topic-name]".
                * /
                "plugin": "status",
                "config": [
                    {
                        /**
                        * How often in minutes to execute command and send result.
                        * /
                        "interval": 10,
                        /**
                        * Set topic where we should publish results.
                        * Use "exequte/system/status/temp" in this example.
                        * /
                        "topic": "temp",
                        /**
                        * Command to execute and return its output.
                        * /
                        "process": "/usr/local/mqtt/reports.sh"
                    },
                    {
                        ... you can specify more actions ...
                    }
                ]
            }
        ]
    }

## TODO

These tasks are of no particular priority:

* provide actual binary packages after code stabilization
* implement HTTP server to ease the configuration process