# Configuration

Supergiant has good defaults and requires only minimal configuration. When in production, however, it may be important to have a deeper understanding of the configuration options we discuss in this document.

## Configuration File

Supergiant ships with an example configuration `config.json.example` file that includes all defaults and that can be downloaded from Supergiant GitHub repository:

```sh
sudo curl https://raw.githubusercontent.com/supergiant/supergiant/master/config/config.json.example --create-dirs -o /etc/supergiant/config.json
```

The file contains host, log, and database settings as well as the list of node sizes for the supported cloud providers (GCE, Amazon AWS, Packet, and DigitalOcean).

## Configuration File Format

Configuration file format is JSON. Here is an example of changing paths to log and database directories using this format:

```json
{
  "sqlite_file": "/var/lib/supergiant/development.db",
  "log_file": "/var/log/supergiant/development.log"
}
```


## Configuration Options

### sqlite_file

Supergiant stores user data like cloud accounts, user information, session, and cluster-related data in the [SQLite](https://www.sqlite.org/about.html) database, a server-less SQL database contained in a single disc file (e.g., `database.db`). The default location for the file is `tmp/development.db`

```json
"sqlite_file": "tmp/development.db" // Type: String
```

### ui_enabled

The parameter defines whether Supergiant UI is enabled. It defaults to `true`

```json
"ui_enabled": "true" // Type: Boolean
```

### capacity_service_enabled

The parameter enables/disables the Supergiant Capacity Service, which automatically creates nodes using CSPs' API depending on RAM and CPU requirements. Turn off the capacity service it's necessary to manage server provisioning manually.

```json
"capacity_service_enabled":"true" // Type: Boolean
```

### publish_host and http_port

The `publish_host` parameter defines the name of a host Supergiant will run on. It defaults to "`localhost`"

```json
"publish_host": "localhost" // Type: String
```

The `http_port` parameter defines the port on which to access Supergiant server running on the published host. It defaults to `8080`

```json
"http_port": "8080" // Type: Integer
```

### log_file

The parameter defines the location of Supergiant log file, and it defaults to `tmp/development.log`

```json
"log_file": "tmp/development.log"
```

### log_level

Supergiant uses a structured logger for Go (golang) [Logrus](https://github.com/Sirupsen/logrus) for logging. The default `log_level` value is `"debug"`.

```json
"log_level": "debug" // Type: String
```

#### Examples of log levels supported by Supergiant:

`debug`: logs information events the most useful to debug the server.

`info`: logs events that describe the progress of user and server tasks.

`error`: logs events that cause errors but still allow the server to run.

### node_sizes

Node sizes parameter defines the list of instance types supported by cloud providers and is used by Capacity Service for auto-scaling and deployment of nodes. See the full list of instance types supported by [AWS](https://aws.amazon.com/ec2/instance-types/), [Packet](https://www.packet.net/bare-metal/), [GCE](https://cloud.google.com/compute/docs/machine-types), and [DigitalOcean](https://www.digitalocean.com/pricing/) in the documentation of these cloud providers.

```json
"node_sizes": {
  "digitalocean": [
      {"name": "s-1vcpu-1gb", "ram_gib": 1, "cpu_cores": 1},
      {"name": "s-1vcpu-2gb", "ram_gib": 2, "cpu_cores": 1},
      {"name": "s-1vcpu-3gb", "ram_gib": 3, "cpu_cores": 1},
      {"name": "s-2vcpu-2gb", "ram_gib": 2, "cpu_cores": 2},
      {"name": "s-3vcpu-1gb", "ram_gib": 1, "cpu_cores": 3},
      {"name": "s-2vcpu-4gb", "ram_gib": 4, "cpu_cores": 2},
      {"name": "s-4vcpu-8gb", "ram_gib": 8, "cpu_cores": 4},
      {"name": "s-6vcpu-16gb", "ram_gib": 16, "cpu_cores": 6},
      {"name": "s-8vcpu-32gb", "ram_gib": 32, "cpu_cores": 8},
      {"name": "s-12vcpu-48gb", "ram_gib": 48, "cpu_cores": 12},
      {"name": "s-16vcpu-64gb", "ram_gib": 64, "cpu_cores": 16},
      {"name": "s-20vcpu-96gb", "ram_gib": 96, "cpu_cores": 20},
      {"name": "s-24vcpu-128gb", "ram_gib": 128, "cpu_cores": 24},
      {"name": "s-32vcpu-192gb", "ram_gib": 192, "cpu_cores": 32}
    ]
}
```

Supergiant's Capacity Service works with the cloud providers APIs to automatically deploy nodes according to RAM and CPU requirements specified when deploying clusters. If Capacity Service is enabled, there's no need to worry about virtual instance flavors used by Supergiant to provision resources.