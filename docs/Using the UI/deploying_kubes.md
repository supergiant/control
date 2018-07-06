# Deploying Kubes

Once your cloud account is linked to Supergiant, you can deploy and auto scale your first Kubernetes cluster (Kube). Supergiant dramatically simplifies this process, taking upon itself all the tedious work of configuring the cluster. All you need to do is to select a cloud account, make some easy edits, and watch Supergiant go. This guide will give you a comprehensive understanding of how to deploy and auto scale Kubernetes clusters with all supported cloud providers.

## Select a Cloud Provider

To begin deploying a kube, first click '**Launch a Cluster**' or  '**Take me to Clusters**'  from your Supergiant's home page. 

![](https://s3-ap-southeast-2.amazonaws.com/sg-github-wiki-images/wiki-deploying-cluster/create-cluster-start.png)

If you haven't deployed any clusters yet, you'll see a '**Deploy your First Cluster**' button. Clicking on it will send you to the page with a list of available cloud accounts. Select a cloud account you want to deploy a Kubernetes cluster to. 

![](https://s3-ap-southeast-2.amazonaws.com/sg-github-wiki-images/wiki-deploying-cluster/deploy-cluster-select-cloud-acc.png)

### Amazon AWS

To deploy a cluster to Amazon AWS, you'll need to configure some parameters for master and node sizes, secure connection, deployment regions, etc. Supergiant will pre-populate some values for you, but be sure to edit them according to your requirements.

![](https://s3-ap-southeast-2.amazonaws.com/sg-github-wiki-images/wiki-deploying-cluster/kube-deploy_aws.gif)

AWS cluster configuration includes the following *parameters*: 

- `Name`  --  a user-friendly name for your Kube. The field accepts a string that starts with a letter, up to 12 characters made of lowercase letters and numbers, and/or dashes `(-)`. 

  **Note**. Dashes are not recommended on AWS because some AWS services might not accept them.

- `Cloud Account Name`  -- a name for your AWS cloud account.

- `Region` -- An AWS region in which to deploy your cluster. The default region is `us-east-1`.  Discover other available regions in the [AWS documentation](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html).

- `VPC IP Range` -- a range of IPv4 addresses for your Virtual Private Cloud (VPC) in a form of a Classless Inter-Domain Routing (CIDR) block. The default range is `172.20.0.0/16`

- `SSH Public Key`  -- an SSH public key that will be used to connect to your cluster. 

- `Master Node Size`  -- the size of the server on which the master will live. Defaults to `m4.large`

- `Kube Master Count`  -- the number of master nodes that will manage your cluster. Defaults to `1`

- `Kubernetes version`:  -- Kubernetes version used for your cluster deployment.  Available options include `1.5.7`, `1.6.7`, `1.7.7` , and `1.8.7`. The default version is `1.5.7`.

- `Node sizes` -- the list of AWS EC2 instance types used by Supergiant to provision new hardware while auto scaling. The list starts with the smallest instance type and ends with the largest. The list is pre-populated with `m4` family that is very reliable running Kubernetes.  See [AWS Instance Types](https://aws.amazon.com/ec2/instance-types/) for a complete list.

### Google Compute Engine (GCE)

GCE configuration for a cluster is similar to other providers with some differences for node types and availability zones. 

The default availability zone for GCE is `us-east1-b`. Consult [GCE documentation](https://cloud.google.com/compute/docs/regions-zones/) for a list of other zones. The default array of node sizes includes `n1-standard` family, which is extremely reliable and cost effective for Kubernetes clusters. The last digit of the instance refers to the number of cores. For example, `n1-standard-8` is a standard machine type with 8 virtual CPUs and 30 GB of RAM. You can examine other [available GCE machine types](https://cloud.google.com/compute/docs/machine-types) in their documentation. 

```json
{
  "cloud_account_name": "gce",
  "gce_config": {
    "ssh_pub_key": "YOUR SSH KEY",
    "zone": "us-east1-b"
  },
  "master_node_size": "n1-standard-1",
  "name": "gcetut",
  "node_sizes": [
    "n1-standard-1",
    "n1-standard-2",
    "n1-standard-4",
    "n1-standard-8"
  ]
}
```

After you submit the form, Supergiant will start provisioning your GCE cluster, a process that usually takes up to 5 minutes. When the process is completed, you'll see new Kubernetes master and minion instances in your GCE console. After that, you can retrieve information about the created cluster from the terminal running the following command.

```sh
curl --insecure https://USERNAME:PASSWORD@MASTER_PUBLIC_IP/api/v1/

```

You should get a JSON response describing the details of your new GCE cluster:

```json
    "username": "YOUR USERNAME",
    "password": "YOUR PASSWORD",
    "heapster_version": "v1.4.0",
    "heapster_metric_resolution": "20s",
    "gce_config": {
        "zone": "us-east1-b",
        "master_instance_group": "https://www.googleapis.com/compute/v1/projects/argon-producer-704/zones/us-east1-b/operations/operation-1520513482047-566e621f93e19-4176219b-118ddc9a",
        "minion_instance_group": "https://www.googleapis.com/compute/v1/projects/argon-producer-724/zones/us-east1-b/operations/operation-1520513482872-566e62205d4c0-a588823c-461e2288",
        "master_nodes": [
            "https://www.googleapis.com/compute/v1/projects/argon-producer-704/zones/us-east1-b/instances/gcetut-master-avkge"
        ],
        "master_name": "gcetut-master-avkge",
        "kube_master_count": 1,
        "ssh_pub_key": "SSH Public Key",
        "kubernetes_version": "1.5.1",
        "etcd_discovery_url": "https://discovery.etcd.io/9ef553b5d21b4ec8f7647667918e40dc",
        "master_private_ip": "10.144.0.3"
    },
    "master_public_ip": "104.196.192.25"
```

### Packet.net 

Packet cluster deployment configuration includes two **sections**: general cluster configuration and Packet-specific settings. 

![](https://s3-ap-southeast-2.amazonaws.com/sg-github-wiki-images/wiki-deploying-cluster/packet-deploy.gif)

#### General Configuration

- `Master Node Size` -- the default value is **Type 0** Packet servers -- [cost-effective bare metal](https://www.packet.net/blog/tiny-meets-mighty-in-our-new-type-0/) servers excellent for testing Kubernetes clusters.

- `Kube Master Count` -- a number of Kube masters. Defaults to `1`. 

- `SSH Pub Key`  --  SSH public key used to connect to the cluster. Packet configuration supports one SSH pub key. 

- `Node Sizes` --  default node sizes for auto-scaling your Packet clusters. Default values are `Type 0`, `Type 1`, `Type 2`, `Type 3`, and `Type 2a`. For more information, see [a full list](https://www.packet.net/bare-metal/servers/tiny/) of Packet bare-metal servers.

#### Packet Configuration

`Facility` -- A Packet region in which the Kube will be created. `ewr1` is the default region. Check out their [official documentation](https://www.packet.net/locations/) for a full list and description of available Packet regions. 

`Project` -- the Packet Project in which the Kube will be created. The field expects a Project UUID similar to  `45acaac1-8adg-7a4f-43da-1aba2183da30` . Your project's UUID can be found under **Project Settings** in Packet's control panel. 

 ![](https://s3-ap-southeast-2.amazonaws.com/sg-github-wiki-images/wiki-deploying-cluster/packet-project-id.png)

### **DigitalOcean** 

The setup for DigitalOcean clusters is essentially the same as for Packet.net. You have several pre-populated fields with options specific to DigitalOcean.

![](https://s3-ap-southeast-2.amazonaws.com/sg-github-wiki-images/wiki-deploying-cluster/digital-ocean-deploy.gif)

#### **Configuration** 

`Region` -- DigitalOcean region in which to deploy a cluster. The default region is `nyc1`. For more information, see [a full list](https://status.digitalocean.com/) of regions available for your DigitalOcean clusters.

`SSH Key Fingertip` -- DigitalOcean SSH key fingertip. The configuration supports the assignment of multiple SSH key fingertips. 

`Master node size`  -- default size is 1GB.

`Kube Master Count` -- default master count is `1`. 

`Node Sizes` -- Supergiant abstracts the DigitalOcean infrastructure  so you can specify node sizes in raw values (e.g., 1 GB, 2GB, 4GB, 8GB, and so on). You can add as many node sizes as you wish up to the largest droplet supported by Digital Ocean. Supergiant will map them to corresponding droplet types stored in the `config.json` file. 

```json
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
```

**Note**: Unlike AWS Kubes, DigitalOcean and Packet.net deployments do not support selection of the  Kubernetes version. 

  ​