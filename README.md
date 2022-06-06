# BruteRatel Notifier

BruteRatel Notifier notifies red team members when their BruteRatel badger status
has been updated. It supports both Slack and Email notification profiles by
default, but it's very extensible so new notification profiles can be added
easily.

This Project is modified from: https://github.com/t94j0/gophish-notifier

## Installation

### From Source

```bash
git clone https://github.com/elephacking/br-notifier
cd br-notifier
go build -o br-notifier
```

## Configuration

The configuration path is `/etc/br_notifier/config.yml`. Below is an example config:

```yaml
# Host to listen on. If GoPhish is running on the same host, you can set this to 127.0.0.1
listen_host: 0.0.0.0
# Port to listen on
listen_port: 9999
# Webhook path. The full address will be http://<host>:<ip><webhook_path>. Ex: http://127.0.0.1:9999/webhook
webhook_path: /webhook
# Log level. Log levels are panic, fatal, error, warning, info, debug, trace.
log_level: debug
# Enables sending profiles. Options are `slack` and `email`. Make sure to configure the required parameters for each profile
profiles:
  - slack
  - email

# Slack Profile
slack:
  # Webhook address. Typically starts with https://hooks.slack.com/services/...
  webhook: '<Your Slack Webhook>'
  # Username displayed in Slack
  username: BRBot
  # Channel to post in
  bot_channel: '#testing2'
  # Bot profile picture
  emoji: ':blowfish:'
  # (Optional) Disable email, username, and credentials from being sent to Slack
  disable_credentials: true

# Email Profile
# Email to send from
email:
  sender: test@test.com
  # Password of sender email. Uses plain SMTP authentication
  sender_password: password123
  # Recipient of notifications
  recipient: mail@example.com
  # Email host to send to
  host: smtp.gmail.com
  # Email host address
  host_addr: smtp.gmail.com:587

# You can also supply an email template for each notification
email_config_template: |
  You caught a new badger!
Badger ID - {{ .Badger }}
Badger User ID - {{ .Config.User_id }}
Badger Hostname - {{ .Config.Hostname }}
Badger Localip - {{ .Config.Localip }}
Badger Process Name - {{ .Config.Process_name }}
Badger Process ID - {{ .Config.Process_id }}
Badger Last Seen - {{ .Config.Last_seen }}
Badger Windows Version - {{ .Config.Windows_version }}
Badger OS Build - {{ .Config.Bld }}
```

Project inspired by [gophish-notifier]

[gophish-notifier]: https://github.com/t94j0/gophish-notifier
[gophish-notifications]: https://github.com/dunderhay/gophish-notifications
