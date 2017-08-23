export class KubesModel {
  aws = {
    'model': {
      'aws_config': {
        'region': 'us-east-1',
        'vpc_ip_range': '172.20.0.0/16'
      },
      'cloud_account_name': '',
      'master_node_size': 'm4.large',
      'name': '',
      'node_sizes': [
        'm4.large',
        'm4.xlarge',
        'm4.2xlarge',
        'm4.4xlarge'
      ]
    },
    'schema': {
      'properties': {
        'aws_config': {
          'properties': {
            'region': {
              'default': 'us-east-1',
              'description': 'Region',
              'type': 'string'
            },
            'vpc_ip_range': {
              'default': '172.20.0.0/16',
              'description': 'VPC IP Range',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'cloud_account_name': {
          'description': 'Cloud Account Name',
          'type': 'string'
        },
        'master_node_size': {
          'default': 'm4.large',
          'description': 'Master Node Size',
          'type': 'string'
        },
        'name': {
          'description': 'Name',
          'type': 'string'
        },
        'node_sizes': {
          'description': 'Node Sizes',
          'type': 'string'
        }
      }
    }
  };

  digitalocean = {
    'model': {
      'cloud_account_name': '',
      'digitalocean_config': {
        'region': 'nyc1',
        'ssh_key_fingerprint': []
      },
      'master_node_size': '1gb',
      'name': '',
      'node_sizes': [
        '1gb',
        '2gb',
        '4gb',
        '8gb',
        '16gb',
        '32gb',
        '48gb',
        '64gb'
      ]
    },
    'schema': {
      'properties': {
        'cloud_account_name': {
          'description': 'Cloud Account Name',
          'type': 'string'
        },
        'digitalocean_config': {
          'properties': {
            'region': {
              'default': 'nyc1',
              'description': 'Region',
              'type': 'string'
            },
            'ssh_key_fingerprint': {
              'description': 'SSH Key Fingerprint',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'master_node_size': {
          'default': '1gb',
          'description': 'Master Node Size',
          'type': 'string'
        },
        'name': {
          'description': 'Name',
          'type': 'string'
        },
        'node_sizes': {
          'description': 'Node Sizes',
          'type': 'string'
        }
      }
    }
  };

  gce = {
    'model': {
      'cloud_account_name': '',
      'gce_config': {
        'ssh_pub_key': '',
        'zone': 'us-east1-b'
      },
      'master_node_size': 'n1-standard-1',
      'name': '',
      'node_sizes': [
        'n1-standard-1',
        'n1-standard-2',
        'n1-standard-4',
        'n1-standard-8'
      ]
    },
    'schema': {
      'properties': {
        'cloud_account_name': {
          'description': 'Cloud Account Name',
          'type': 'string'
        },
        'gce_config': {
          'properties': {
            'ssh_pub_key': {
              'description': 'SSH Public Key',
              'type': 'string'
            },
            'zone': {
              'default': 'us-east1-b',
              'description': 'Zone',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'master_node_size': {
          'default': 'n1-standard-1',
          'description': 'Master Node Size',
          'type': 'string'
        },
        'name': {
          'description': 'Name',
          'type': 'string'
        },
        'node_sizes': {
          'description': 'Node Sizes',
          'type': 'string'
        }
      }
    }
  };

  openstack = {
    'model': {
      'cloud_account_name': '',
      'master_node_size': 'm1.smaller',
      'name': '',
      'node_sizes': [
        'm1.smaller',
        'm1.small'
      ],
      'openstack_config': {
        'image_name': 'CoreOS',
        'region': 'RegionOne',
        'ssh_key_fingerprint': ''
      },
      'ssh_pub_key': ''
    },
    'schema': {
      'properties': {
        'cloud_account_name': {
          'description': 'Cloud Account Name',
          'type': 'string'
        },
        'master_node_size': {
          'default': 'm1.smaller',
          'description': 'Master Node Size',
          'type': 'string'
        },
        'name': {
          'description': 'Name',
          'type': 'string'
        },
        'node_sizes': {
          'description': 'Node Sizes',
          'type': 'string'
        },
        'openstack_config': {
          'properties': {
            'image_name': {
              'default': 'CoreOS',
              'description': 'Image Name',
              'type': 'string'
            },
            'region': {
              'default': 'RegionOne',
              'description': 'Region',
              'type': 'string'
            },
            'ssh_key_fingerprint': {
              'description': 'SSH Key Fingerprint',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'ssh_pub_key': {
          'description': 'SSH Public Key',
          'type': 'string'
        }
      }
    }
  };

  packet = {
    'model': {
      'cloud_account_name': '',
      'master_node_size': 'Type 0',
      'name': '',
      'node_sizes': [
        'Type 0',
        'Type 1',
        'Type 2',
        'Type 3'
      ],
      'packet_config': {
        'facility': 'ewr1',
        'kube_master_count': 1,
        'project': '',
        'ssh_pub_key': ''
      }
    },
    'schema': {
      'properties': {
        'cloud_account_name': {
          'description': 'Cloud Account Name',
          'type': 'string'
        },
        'master_node_size': {
          'default': 'Type 0',
          'description': 'Master Node Size',
          'type': 'string'
        },
        'name': {
          'description': 'Name',
          'type': 'string'
        },
        'node_sizes': {
          'description': 'Node Sizes',
          'type': 'string'
        },
        'packet_config': {
          'properties': {
            'facility': {
              'default': 'ewr1',
              'description': 'Facility',
              'type': 'string'
            },
            'kube_master_count': {
              'description': 'Kube Master Count',
              'type': 'string'
            },
            'project': {
              'description': 'Project',
              'type': 'string'
            },
            'ssh_pub_key': {
              'description': 'SSH Public Key',
              'type': 'string'
            }
          },
          'type': 'object'
        }
      }
    }
  };
  public providers = {
    'aws': this.aws,
    'digitalocean': this.digitalocean,
    'gce': this.gce,
    'openstack': this.openstack,
    'packet': this.packet
  };
}
