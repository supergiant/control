# Linux

1. Download the appropriate Supergiant API and UI binaries from the [tags page](https://github.com/supergiant/supergiant/releases) (until our upcoming 1.0 release, only Linux is available). For example: 

```sh
curl https://github.com/supergiant/supergiant/releases/download/v1.0.0-beta.4/supergiant-api-linux-amd64 -L -o /usr/bin/supergiant-api
curl https://github.com/supergiant/supergiant/releases/download/v1.0.0-beta.4/supergiant-ui-linux-amd64 -L -o /usr/bin/supergiant-ui
```

2. Grant executability to the binaries:

```sh
sudo chmod +x /usr/bin/supergiant-api /usr/bin/supergiant-ui
```

3. Download the [example configuration file](https://github.com/supergiant/supergiant/blob/master/config/config.json.example): 

```sh
curl https://raw.githubusercontent.com/supergiant/supergiant/master/config/config.json.example --create-dirs -o /etc/supergiant/config.json
```

4. In the configuration file, specify log and database paths (for more detailed information configuration, see [the detailed doc](http://supergiant.readthedocs.io/en/docs/Installation/Configuration/)):

```json
{
 ...
 "sqlite_file": "/var/lib/supergiant/development.db",
 ...
 "log_file": "/var/log/supergiant/development.log",
 ...
}
```

5. Run the API with a config file and save the generated username and password (it will only be shown once):

```sh
/usr/bin/supergiant-api --config-file /etc/supergiant/config.json
```

6. Tail the log in another window or session to see that Supergiant is running well:
```sh
tail -f /var/log/supergiant/development.log
```

7. Run the UI in another window or session:
```sh
/usr/bin/supergiant-ui
```

8. Access Supergiant's UI on port 3001!
