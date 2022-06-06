package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var defaultConfigTemplate = `Badger ID - {{ .Badger }}
Badger OS Build - {{ .Config.Bld }}
Badger C2 - {{ .Config.C2 }}
Badger C2 ID - {{ .Config.C2_id }}
Badger Cookie - {{ .Config.Cookie }}
Badger Hostname - {{ .Config.Hostname }}
Badger Localip - {{ .Config.Localip }}
Badger Process Name - {{ .Config.Process_name }}
Badger Process ID - {{ .Config.Process_id }}
Badger Last Seen - {{ .Config.Last_seen }}
Badger User ID - {{ .Config.User_id }}
Badger Windows Version - {{ .Config.Windows_version }}
Badger is Dead? - {{ .Config.Dead }}
Badger Is Pvt?", {{ .Config.Is_pvt }}
Badger Pipeline - {{ .Config.Pipeline }}
Badger Pvt Master - {{ .Config.Pvt_master }}`

var defaultMessageTemplate = `Email ID - {{ .ID }}
Email Address - {{ .Email }}
IP Address - {{ .Address }}
User Agent - {{ .UserAgent }}
Username - {{ .Username }}
Password - {{ .Password }}`

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/br_notifier")
	viper.AddConfigPath(".")
	setDefaults()
	if err := viper.ReadInConfig(); err != nil {
		log.Error(err)
	}
	log.Infof("Using config file: %s", viper.ConfigFileUsed())
	validateConfig()
	setLogLevel()
}

func setDefaults() {
	viper.SetDefault("log_level", "info")
	viper.SetDefault("slack.bot_username", "BRBot")
	viper.SetDefault("slack.bot_emoji", ":blowfish:")
	viper.SetDefault("listen_host", "0.0.0.0")
	viper.SetDefault("listen_port", "9999")
	viper.SetDefault("webhook_path", "/webhook")
	viper.SetDefault("email_config_template", defaultConfigTemplate)
	viper.SetDefault("email_message_template", defaultMessageTemplate)
	viper.SetDefault("profiles", []string{"slack"})
}

func setLogLevel() {
	level, err := log.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		log.Fatal("log level must be a valid level: panic, fatal, error, warning, info, debug, trace")
	}
	log.SetLevel(level)
}

func validateConfig() {
	checkKeysExist := func(keys ...string) {
		for _, key := range keys {
			if !viper.IsSet(key) {
				log.Fatal("Config value is not set: ", key)
			}
		}
	}

	globalConfigs := []string{"profiles"}
	checkKeysExist(globalConfigs...)

	profiles := viper.GetStringSlice("profiles")
	for _, profile := range profiles {
		if profile == "slack" {
			slackConfigs := []string{"slack.webhook", "slack.bot_channel"}
			checkKeysExist(slackConfigs...)
			log.Infof("Using Slack sending profile. Will send messages to %s", viper.GetString("slack.bot_channel"))
			continue
		}
		if profile == "email" {
			emailConfigs := []string{"email.sender", "email.sender_password", "email.recipient", "email.host", "email.host_addr"}
			checkKeysExist(emailConfigs...)
			log.Infof("Using Email sending profile. Will send emails from %s to %s",
				viper.GetString("email.sender"),
				viper.GetString("email.recipient"))
			continue
		}
		log.Fatalf("Profile \"%s\" does not exist", profile)
	}
}
