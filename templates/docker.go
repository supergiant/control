package templates

const dockerTpl = `
DOCKER_VERSION={{ .Version }}
ARCH={{ .Arch }}

sudo apt-get update -y
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo apt-key fingerprint 0EBFCD88

sudo add-apt-repository \
	"deb [arch=${ARCH}] https://download.docker.com/linux/ubuntu \
	$(lsb_release -cs) \
	stable"

sudo apt-get update -y

# show available docker versions:
# apt-cache madison docker-ce

FULL_DOCKER_VERSION=$(apt-cache madison docker-ce | cut -d '|' -f2 | tr -d ' ' | grep "${DOCKER_VERSION}")
if [ -z "${FULL_DOCKER_VERSION}" ]; then
	echo "package for the ${DOCKER_VERSION} docker version not found"
	echo "Available packages:"
	apt-cache madison docker-ce | cut -d '|' -f2 | tr -d ' '
	exit 1
fi

sudo apt-get install -y docker-ce=${FULL_DOCKER_VERSION} containerd.io
`
