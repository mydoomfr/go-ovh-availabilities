# go-ovh-availabilities

## Summary
**go-ovh-availabilities** is a cross-platform monitoring tool which aims to alert you from OVH dedicated new availabilities. You can be alerted by mail or simply from the terminal log output.


![screenshot.png](screenshot.png)

```
$ go-ovh-availabilities -c /path/to/config.json
2018/10/09 20:29:17 Configuration file : config.json
2018/10/09 20:29:18 KS-3 1801sk14 is now : available (1H-low) in europe (gra) !
2018/10/09 20:29:18 KS-2 1801sk13 is now : available (1H-low) in europe (rbx) !
2018/10/09 20:29:18 1801sys45 is now : available (240H) in europe (rbx) !
```

## Requirements :

Mandatory :

* Go >= 1.10

Optional : 

* SMTP server (for mail support)
* Sendmail (for easier mail support) - [Linux-only]


## Usage

```shell
$ go install github.com/mydoom/go-ovh-availabilities
$ vi config.json
$ go-ovh-availabilities -c /path/to/config.json
```

Well tested on : W10, OSX Mojave & Ubuntu 16.04

## Configuration

- Refresh (mandatory) is the interval to wait before the next API request. (seconds) Warning: Do not use low-value or you may be automatically banned.
- Wanted.Name (optional) is your custom search name.
- Wanted.Hardware (mandatory) is the "id" of the dedicated server. You can get it from the Kimsufi or Souyoustart website by inspecting the source code.
- Wanted.Region (optional) refers to the geolocation.
- Wanted.datacenters (optional) is useful if you want a desired DC. Example: "rbx" for "Roubaix"...

### config.json 

```json
{
  "api": "https://www.ovh.com/engine/api/dedicated/server/availabilities?country=fr",
  "refresh": 1800,
  "wanted": [
    {
      "name": "My favorite KS-1",
      "hardware": "1801sk12"
    },
    {
      "name": "Backup storage (KS-2)",
      "hardware": "1801sk13",
      "region": "europe",
      "datacenters": ["rbx", "lon"]
    },
    {
      "hardware": "1801sys011",
      "region": "europe",
      "datacenters": []
    }
  ]
}
```

## Mail Support

### Method #1 : Sendmail

Linux-only. May work with OSX if you have sendmail.

Add to config.json :

```json
{
  "mail": {
    "from": "user@domain.tld",
    "to": "someone@exemple.com",
    "object": "Alert : OVH",
    "sendmail": {
      "active": true,
      "bin": "/usr/sbin/sendmail"
    }
  }
}
```

### Method #2 : SMTP Auth

Add to config.json :

```json
{
  "mail": {
    "from": "user@domain.tld",
    "to": "someone@exemple.com",
    "object": "Alert : OVH",
    "smtp": {
      "active": true,
      "server": "localhost",
      "port": 1025,
      "username": "",
      "password": ""
    }
  }
}
```