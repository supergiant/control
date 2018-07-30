export class ClusterGCEModel {
  gce = {
    'data': {
      'cloud_account_name': '',
      'gce_config': {
        'zone': 'us-east1-b'
      },
      'ssh_pub_key': '',
      'master_node_size': 'n1-standard-1',
      'name': '',
      'kube_master_count': 1,
      'kubernetes_version': '1.8.7',
      'node_sizes': [
        'n1-standard-1',
        'n1-standard-2',
        'n1-standard-4',
        'n1-standard-8'
      ]
    },
    'schema': {
      'properties': {
        'name': {
          'type': 'string',
          'description': 'The desired name of the kube. Max length of 12 characters.',
          'pattern': '^[a-z]([-a-z0-9]*[a-z0-9])?$',
          'maxLength': 12
        },
        'kubernetes_version': {
          'type': 'string',
          'description': 'The Version of Kubernetes to be deployed.',
          'default': '1.5.7',
          'enum': ['1.5.7', '1.6.7', '1.7.7', '1.8.7', '1.9.7', '1.10.6', '1.11.1'] // TODO: <-- Should be dynamically populated.
        },
        'cloud_account_name': {
          'type': 'string',
          'description': 'The Supergiant cloud account you created for use with GCE.'
        },
        'gce_config': {
          'type': 'object',
          'properties': {
            'zone': {
              'type': 'string',
              'description': 'The GCE zone the kube will be created in.',
              'default': 'us-east1-b',
              'oneOf': [
                // northamerica-northeast1
                { 'title': 'northamerica-northeast1: northamerica-northeast1-a', 'const': 'northamerica-northeast1-a' },
                { 'title': 'northamerica-northeast1: northamerica-northeast1-b', 'const': 'northamerica-northeast1-b' },
                { 'title': 'northamerica-northeast1: northamerica-northeast1-c', 'const': 'northamerica-northeast1-c' },
                // us-west1
                { 'title': 'us-west1: us-west1-a', 'const': 'us-west1-a' },
                { 'title': 'us-west1: us-west1-b', 'const': 'us-west1-b' },
                { 'title': 'us-west1: us-west1-c', 'const': 'us-west1-c' },
                // us-central1
                { 'title': 'us-central1: us-central1-a', 'const': 'us-central1-a' },
                { 'title': 'us-central1: us-central1-b', 'const': 'us-central1-b' },
                { 'title': 'us-central1: us-central1-c', 'const': 'us-central1-c' },
                { 'title': 'us-central1: us-central1-f', 'const': 'us-central1-f' },
                // us-east1
                { 'title': 'us-east1: us-east1-b', 'const': 'us-east1-b' },
                { 'title': 'us-east1: us-east1-c', 'const': 'us-east1-c' },
                { 'title': 'us-east1: us-east1-d', 'const': 'us-east1-d' },
                // us-east4
                { 'title': 'us-east4: us-east4-a', 'const': 'us-east4-a' },
                { 'title': 'us-east4: us-east4-b', 'const': 'us-east4-b' },
                { 'title': 'us-east4: us-east4-c', 'const': 'us-east4-c' },
                // southamerica-east1
                { 'title': 'southamerica-east1: southamerica-east1-a', 'const': 'southamerica-east1-a' },
                { 'title': 'southamerica-east1: southamerica-east1-b', 'const': 'southamerica-east1-b' },
                { 'title': 'southamerica-east1: southamerica-east1-c', 'const': 'southamerica-east1-c' },
                // europe-north1
                { 'title': 'europe-north1: europe-north1-a', 'const': 'europe-north1-a' },
                { 'title': 'europe-north1: europe-north1-b', 'const': 'europe-north1-b' },
                { 'title': 'europe-north1: europe-north1-c', 'const': 'europe-north1-c' },
                // europe-west1
                { 'title': 'europe-west1: europe-west1-b', 'const': 'europe-west1-b' },
                { 'title': 'europe-west1: europe-west1-c', 'const': 'europe-west1-c' },
                { 'title': 'europe-west1: europe-west1-d', 'const': 'europe-west1-d' },
                // europe-west2
                { 'title': 'europe-west2: europe-west2-a', 'const': 'europe-west2-a' },
                { 'title': 'europe-west2: europe-west2-b', 'const': 'europe-west2-b' },
                { 'title': 'europe-west2: europe-west2-c', 'const': 'europe-west2-c' },
                // europe-west3
                { 'title': 'europe-west3: europe-west3-a', 'const': 'europe-west3-a' },
                { 'title': 'europe-west3: europe-west3-b', 'const': 'europe-west3-b' },
                { 'title': 'europe-west3: europe-west3-c', 'const': 'europe-west3-c' },
                // europe-west4
                { 'title': 'europe-west4: europe-west4-a', 'const': 'europe-west4-a' },
                { 'title': 'europe-west4: europe-west4-b', 'const': 'europe-west4-b' },
                { 'title': 'europe-west4: europe-west4-c', 'const': 'europe-west4-c' },
                // asia-northeast1
                { 'title': 'asia-northeast1: asia-northeast1-a', 'const': 'asia-northeast1-a' },
                { 'title': 'asia-northeast1: asia-northeast1-b', 'const': 'asia-northeast1-b' },
                { 'title': 'asia-northeast1: asia-northeast1-c', 'const': 'asia-northeast1-c' },
                // asia-east1
                { 'title': 'asia-east1: asia-east1-a', 'const': 'asia-east1-a' },
                { 'title': 'asia-east1: asia-east1-b', 'const': 'asia-east1-b' },
                { 'title': 'asia-east1: asia-east1-c', 'const': 'asia-east1-c' },
                // asia-southeast1
                { 'title': 'asia-southeast1: asia-southeast1-a', 'const': 'asia-southeast1-a' },
                { 'title': 'asia-southeast1: asia-southeast1-b', 'const': 'asia-southeast1-b' },
                { 'title': 'asia-southeast1: asia-southeast1-c', 'const': 'asia-southeast1-c' },
                // asia-south1
                { 'title': 'asia-south1: asia-south1-a', 'const': 'asia-south1-a' },
                { 'title': 'asia-south1: asia-south1-b', 'const': 'asia-south1-b' },
                { 'title': 'asia-south1: asia-south1-c', 'const': 'asia-south1-c' },
                // australia-southeast1
                { 'title': 'australia-southeast1: australia-southeast1-a', 'const': 'australia-southeast1-a' },
                { 'title': 'australia-southeast1: australia-southeast1-b', 'const': 'australia-southeast1-b' },
                { 'title': 'australia-southeast1: australia-southeast1-c', 'const': 'australia-southeast1-c' },
              ] // TODO: if we keep using this tool, split into two selects - region and zone (zone populated by region)
            }
          },
          'required': ['zone']
        },
        'ssh_pub_key': {
          'type': 'string',
          'description': 'The public key that will be used to SSH into the kube.',
          'widget': 'textarea'
        },
        'master_node_size': {
          'type': 'string',
          'default': 'n1-standard-1',
          'description': 'The size of the server the master will live on.'
        },
        'kube_master_count': {
          'type': 'number',
          'description': 'The number of masters desired - for High Availability.',
          'widget': 'number',
          'minimum': 1
        },
        'node_sizes': {
          'type': 'array',
          'description': 'The sizes you want to be available to Supergiant when scaling.',
          'id': '/properties/node_sizes',
          'items': {
            'id': '/properties/node_sizes/items',
            'type': 'string'
          }
        }
      },
      'required': ['name']
    },
    'layout': [
      {
        'type': 'section',
        'title': 'Cluster Details',
        'items': [
          { 'key': 'name' },
          { 'key': 'kubernetes_version' },
          // { 'key': 'cloud_account_name', 'readonly': true }
        ]
      },
      {
        'type': 'section',
        'title': 'GCE Config',
        'items': [
          { 'key': 'gce_config.zone' },
        ]
      },
      {
        'type': 'section',
        'items': [
          { 'key': 'ssh_pub_key' },
          { 'key': 'master_node_size' },
          { 'key': 'kube_master_count' },
        ]
      },
      {
        'key': 'node_sizes',
        'type': 'array',
        'items': [
          { 'key': 'node_sizes[]' },
        ]
      },
      {
        'type': 'submit',
        'title': 'Create'
      }
    ]
  };
}
