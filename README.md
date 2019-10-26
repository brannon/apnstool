![](https://github.com/brannon/apnstool/workflows/ci-build/badge.svg)

#  Command-line tool for interacting with APNs

`apnstool` is a command-line tool for interacting with APNs (the Apple Push Notification service). It makes it simple to test your APNs credentials and different notification types.

## Why?

There are several web-based tools for testing APNs. They are fairly limited (don't support headers, don't support token-based credentials) and require you to provide your APNs credentials to a 3rd-party website.

This tool is open-source and meant to be run from your computer. It supports most common features of APNs (and can be extended to support more).

## This tool is not meant for production usage

There are many production-ready options for interacting with APNs that offer cross-platform notifications, advanced device management, and multi-targeting.

I'm partial to [Azure Notification Hubs](https://azure.microsoft.com/en-us/services/notification-hubs/).

## TODO

- Add support for APNs certificate-based authentication.
- Add support for saving credentials.
- Add support for getting flags from config/environment.
- Add support for reading the `--data` value from stdin.
- Add support for arbitrary `apns-` headers.
