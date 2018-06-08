#!/bin/bash -e
# Initilize variables
init_vars() {
  # For signature v4 signning purpose
  timestamp=$(date -u "+%Y-%m-%d %H:%M:%S")
  isoTimpstamp=$(date -ud "${timestamp}" "+%Y%m%dT%H%M%SZ")
  dateScope=$(date -ud "${timestamp}" "+%Y%m%d")
  #dateHeader=$(date -ud "${timestamp}" "+%a, %d %h %Y %T %Z")
  signedHeaders="host;x-amz-content-sha256;x-amz-date;x-amz-security-token"
  service="s3"

  # Get instance auth token from meta-data
  region=$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/document/ | jq -r '.region')

  # Bucket
  bucket={{ .AWSConfig.BucketName }}
  roleProfile=kubernetes-master
  # KeyId, secret, and token
  accessKeyId=$(curl -s http://169.254.169.254/latest/meta-data/iam/security-credentials/$roleProfile | jq -r '.AccessKeyId')
  secretAccessKey=$(curl -s http://169.254.169.254/latest/meta-data/iam/security-credentials/$roleProfile | jq -r '.SecretAccessKey')
  stsToken=$(curl -s http://169.254.169.254/latest/meta-data/iam/security-credentials/$roleProfile | jq -r '.Token')

  # Path to cloud-config.yaml. e.g. worker/cloud-config.yaml
  cloudConfigYaml="build/master.yaml"

  # Path to initial-cluster urls file to join cluster
  initialCluster="etcd/initial-cluster"

  workDir="/tmp"

  # Empty payload hash (we are getting content, not upload)
  payload=$(sha256_hash /dev/null)


  if [[ $region == "us-east-1" ]]; then
   PREFIX="s3"
  else
   PREFIX="s3-${region}"
  fi
  # Host header
  hostHeader="${bucket}.${PREFIX}.amazonaws.com"

  # Curl options
  opts="-v -L --fail --retry 5 --retry-delay 3 --silent --show-error"
  # Curl logs
  bootstrapLog=${workDir}/bootstrap.log
}

# Untilities
hmac_sha256() {
  key="$1"
  data="$2"
  echo -n "$data" | openssl dgst -sha256  -mac HMAC -macopt "$key" | sed 's/^.* //'
}
sha256_hash() {
  echo $(sha256sum "$1" | awk '{print $1}')
}

curl_get() {
curl $opts -H "Host: ${hostHeader}" \
    -H "Authorization: AWS4-HMAC-SHA256 \
    Credential=${accessKeyId}/${dateScope}/${region}/s3/aws4_request, \
    SignedHeaders=${signedHeaders}, Signature=${signature}" \
    -H "x-amz-content-sha256: ${payload}" \
    -H "x-amz-date: ${isoTimpstamp}" \
    -H "x-amz-security-token:${stsToken}" \
     https://${hostHeader}/${filePath}
}

canonical_request() {
  echo "GET"
  echo "/${filePath}"
  echo ""
  echo host:${hostHeader}
  echo "x-amz-content-sha256:${payload}"
  echo "x-amz-date:${isoTimpstamp}"
  echo "x-amz-security-token:${stsToken}"
  echo ""
  echo "${signedHeaders}"
  printf "${payload}"
}

string_to_sign() {
  echo "AWS4-HMAC-SHA256"
  echo "${isoTimpstamp}"
  echo "${dateScope}/${region}/s3/aws4_request"
  printf "$(canonical_request | sha256_hash -)"
}

signing_key() {
  dateKey=$(hmac_sha256 key:"AWS4$secretAccessKey" $dateScope)
  dateRegionKey=$(hmac_sha256 hexkey:$dateKey $region)
  dateRegionServiceKey=$(hmac_sha256 hexkey:$dateRegionKey $service)
  signingKey=$(hmac_sha256 hexkey:$dateRegionServiceKey "aws4_request")
  printf "${signingKey}"
}

# Initlize varables
init_vars

cd ${workDir}

## Download File
filePath=${cloudConfigYaml}
signature=$(string_to_sign | openssl dgst -sha256 -mac HMAC -macopt hexkey:$(signing_key) | awk '{print $NF}')
curl_get 2>> ${bootstrapLog} > ${workDir}/cloud-config.yaml
# Run cloud-init
coreos-cloudinit --from-file=/tmp/cloud-config.yaml