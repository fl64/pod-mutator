TEMP_DIR=/tmp/k8s-webhook-server/serving-certs/
mkdir -p ${TEMP_DIR}
cfssl gencert -initca ca-csr.json | cfssljson -bare ${TEMP_DIR}/ca –
cfssl gencert -ca=${TEMP_DIR}/ca.pem -ca-key=${TEMP_DIR}/ca-key.pem -config=ca-config.json -profile=server mwh-csr.json | cfssljson -bare ${TEMP_DIR}/tls
mv ${TEMP_DIR}/tls-key.pem ${TEMP_DIR}/tls.key
mv ${TEMP_DIR}/tls.pem ${TEMP_DIR}/tls.crt
ls ${TEMP_DIR}
