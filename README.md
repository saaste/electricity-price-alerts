# Electricity Price Alerts
Little [Go](https://go.dev/) tool that gets electricity prices from [porssisahko.net](https://porssisahko.net/) and sends notifications if tomorrow's prices are above the price threshold.

I made the script for myself, because I don't want to spend mental energy on actively checking the apps for current prices. So this script is **not** for real-time notifications. Its main purpose is to warn about the expensive periods of the next day.

## Notifications

The following notification is sent for **each** expensive period:

In Finnish:
> Sähkön hinta
>
> Yli hintarajan (12 ¢/kWh) kello 09:00. Alle hintarajan kello 16:00.

In English:
> Electricity price
>
> Above the price threshold (12 ¢/kWh) at 09:00. Below the price threshold at 16:00.

## Requirements
You need to have [Go](https://go.dev/) installed.

For notifications you need to have a [Gotify](https://gotify.net/) server up and running. See the [documentation](https://gotify.net/docs/index) for more information.

## How to run
First, build the binary with
```
go build
```

Then, run it with a scheduling program such as `crontab`. Note that you only need to run it **once a day**. The price data is updated around 14:00 (Europe/Helsinki), so you should schedule it to run daily after 14:00.

```
Usage: electricity-price-alerts --threshold THRESHOLD --gotify GOTIFY --key KEY [--lang LANG]

Options:
  --threshold THRESHOLD, -t THRESHOLD   Price threshold as ¢/kWh
  --gotify GOTIFY, -g GOTIFY            Gotify URL
  --key KEY, -k KEY                     Gotify API key
  --lang LANG, -l LANG                  Notification language [fi, en]. Default: fi
  --help, -h                            display help
```

Example:
```
./electricity-price-alerts -t 12 -g https://mygotify.homeserver.lan -k asFGx3.Zdf3 -l fi
```