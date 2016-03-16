curl -XPOST localhost:8080/registries/dockerhub/repos -d '{
  "name": "qbox",
  "key": "eyJodHRwczovL2luZGV4LmRvY2tlci5pby92MS8iOnsiYXV0aCI6ImMzQmhibXRsYm5OMGFXVnVPbTFsZEdWMU16azUiLCJlbWFpbCI6Im1pa2VAcWJveC5pbyJ9fQ=="
}'

curl -XPOST localhost:8080/apps -d '{
  "name": "test"
}'

curl -XPOST localhost:8080/apps/test/components -d '{
  "name": "elasticsearch",
  "instances": 1,
  "termination_grace_period": 10,
  "volumes": [
    {
      "name": "elasticsearch-data",
      "type": "gp2",
      "size": 20
    }
  ],
  "containers": [
    {
      "image": "qbox/qbox-docker:2.1.1",
      "cpu": {
        "min": 0,
        "max": 1000
      },
      "ram": {
        "min": 3072,
        "max": 4096
      },
      "mounts": [
        {
          "volume": "elasticsearch-data",
          "path": "/data-1"
        }
      ],
      "ports": [
        {
          "protocol": "HTTP",
          "number": 9200,
          "public": true
        },
        {
          "protocol": "TCP",
          "number": 9300
        }
      ],
      "env": [
        {
          "name": "CLUSTER_ID",
          "value": "SG_TEST"
        },
        {
          "name": "NODE_NAME",
          "value": "SG_TEST_instance_id"
        },
        {
          "name": "MASTER_ELIGIBLE",
          "value": "true"
        },
        {
          "name": "DATA_PATHS",
          "value": "/data-1"
        },
        {
          "name": "UNICAST_HOSTS",
          "value": "elasticsearch.test.svc.cluster.local:9300"
        },
        {
          "name": "MIN_MASTER_NODES",
          "value": "1"
        },
        {
          "name": "CORES",
          "value": "1"
        },
        {
          "name": "ES_HEAP_SIZE",
          "value": "2048m"
        },
        {
          "name": "INDEX_NUMBER_OF_SHARDS",
          "value": "4"
        },
        {
          "name": "INDEX_NUMBER_OF_REPLICAS",
          "value": "0"
        }
      ]
    }
  ]
}'
