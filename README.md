```shell
# Run etcd in background (or in another tab/pane)
etcd &

ETCD_ENDPOINT=http://localhost:2379 \
AWS_REGION=us-east-1 \
AWS_AZ=us-east-1c \
AWS_SG_ID=<security_group_id> \
AWS_SUBNET_ID=<subnet_id> \
K8S_HOST=<kube_master_ip> \
K8S_USER=<kube_http_basic_username> \
K8S_PASS=<kube_http_basic_password> \
go run main.go
```

See [example.sh](example.sh) and [api/router.go](api/router.go).


*Note: I'm going to squash commit history soon to get rid of the private
Dockerhub key in example.sh.*

# Server Options

```NAME:
   supergiant-api - The Supergiant api server.

USAGE:
   supergiant [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --etcd-host [--etcd-host option --etcd-host option]	Array of etcd hosts. [$ETCD_ENDPOINT]
   --k8sHost, --kh "kubernetes"				IP of a Kuberntes api. [$K8S_HOST]
   --k8sUser, --ku "<Kubernetes api userID>"		Username used to connect to your Kubernetes api. [$K8S_USER]
   --k8sPass, --kp "<Kubernetes api password>"		Password used to connect to your Kubernetes api. [$K8S_PASS]
   --awsRegion, --ar "<AWS Region>"			AWS Region in which your kubernetes cluster resides. [$AWS_REGION]
   --awsAZ, --az "<AWS Availability Zone>"		AWS Availability Zone in which your kubernetes cluster resides. [$AWS_AZ]
   --help, -h						show help
   --version, -v					print the version
   ```



   #License

   <pre>
   This software is licensed under the Apache License, version 2 ("ALv2"), quoted below.

   Copyright 2016 Qbox, Inc., a Delaware corporation. All rights reserved.

   Licensed under the Apache License, Version 2.0 (the "License"); you may not
   use this file except in compliance with the License. You may obtain a copy of
   the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
   WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
   License for the specific language governing permissions and limitations under
   the License.
   </pre>
