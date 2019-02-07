#!/bin/sh
# https://download.docker.com/linux/ubuntu/dists/xenial/pool/stable/amd64/docker-ce_18.06.1~ce~3-0~ubuntu_amd64.deb

DOCKER_VERSION={{ .Version }}
UBUNTU_RELEASE={{ .ReleaseVersion }}
ARCH={{ .Arch }}
OUT_DIR=/tmp
URL="https://download.docker.com/linux/ubuntu/dists/${UBUNTU_RELEASE}/pool/stable/${ARCH}/docker-ce_${DOCKER_VERSION}~ce-0~ubuntu_${ARCH}.deb"

sudo wget -O $OUT_DIR/$(basename $URL) $URL
sudo apt install -y $OUT_DIR/$(basename $URL)
sudo rm $OUT_DIR/$(basename $URL)

sudo systemctl enable docker.service