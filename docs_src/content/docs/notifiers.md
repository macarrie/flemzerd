+++
title = "Notifiers"
description = "Send alerts for events occuring in flemzerd."
date = 2018-05-24T14:52:30Z
weight = 70
draft = false
bref = "Notifiers define notification ways to use for flemzerd events. Multiple notifiers can be used at once to have multiple platforms notified at the same time."
toc = true
+++

## Notifiers overview
---

Flemzerd sends notifications for events occuring (download start/success/failure, new episodes/movies) using Notifiers. Each Notifier represents a method of notification (mail, sms, Telegram, ...).

Notifiers are optionnal in flemzerd. If no Notifiers are defined, notifications won't be sent. When multiple Notifiers are defined, each notification will be sent using all defined notifiers.

## Available Notifiers
---

### Event logger
---
Flemzerd automatically uses the Event log notifier to make notifications display in the Web interface. This Notifier cannot be disabled.


### Pushbullet
---
 Sends notifications through Pushbullet.

#### How to use
---
* Enable `pushbullet` Notifier in configuration file
{{< highlight toml >}}
[notifiers]
    [notifiers.pushbullet]
        accesstoken = "accesstoken"
{{< /highlight >}}
* The Pushbullet access token can be foudn on your Pushbullet account settings

### Desktop
---
 Sends notifications to the desktop. The way this Notifier sends notifications depends on the OS:
 * Linux: `notify-send` fo Gnome and `kdialog` for KDE
 * Mac OS X: `terminal-notifier` if installed or `osascript` otherwise
 * Windows: `growlnotify`

#### How to use
---
* Enable `desktop` Notifier in configuration file
{{< highlight toml >}}
[notifiers]
    desktop = []
{{< /highlight >}}

### Kodi
---
 Sends notifications to a Kodi instance.

#### How to use
---
* Enable `kodi` Notifier in configuration file
{{< highlight toml >}}
[notifiers]
    kodi = []
{{< /highlight >}}
* Define `kodi` Mediacenter in configuration file. The `kodi` Notifier uses the configuration from the Kodi mediacenter
{{< highlight toml >}}
[mediacenters]
    [mediacenters.kodi]
        address = "address"
        port = 9090
{{< /highlight >}}

### Telegram
---
 Sends notifications via Telegram.

#### How to use
---
* Enable `telegram` Notifier in configuration file
{{< highlight toml >}}
[notifiers]
    telegram = []
{{< /highlight >}}
* Define Telegram bot token
    Flemzerd needs the bot token to be able to send notifications. This token is passed to flemzerd by defining the `FLZ_TELEGRAM_BOT_TOKEN` environment variable. When defined during flemzerd compilation, this variable is compiled into the binary.
    In [packages](https://github.com/macarrie/flemzerd/releases) found on GitHub, this key is precompiled into the binary.
    The Telegram notifier then needs to be enabled in the Settings page of the Web UI.


