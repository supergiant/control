KUBERNETES_SSL_DIR=${KUBERNETES_CONFIG_DIR}/ssl

mkdir -p ${KUBERNETES_SSL_DIR}
echo "{{ .CACert }}" > ${KUBERNETES_SSL_DIR}/'{{ .CACertName }}'
echo "{{ .CAKeyCert }}" > ${KUBERNETES_SSL_DIR}/'{{ .CAKeyName }}'
echo "{{ .APIServerCert }}" > ${KUBERNETES_SSL_DIR}/'{{ .APIServerCertName }}'
echo "{{ .APIServerKey }}" > ${KUBERNETES_SSL_DIR}/'{{ .APIServerKeyName }}'
echo "{{ .KubeletClientCert }}" > ${KUBERNETES_SSL_DIR}/'{{ .KubeletClientCertName }}'
echo "{{ .KubeletClientKey }}" > ${KUBERNETES_SSL_DIR}/'{{ .KubeletClientKeyName }}'