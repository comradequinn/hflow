#! /bin/bash

CA="hflow-ca"
DIR="./bin/certs/ca"

rm -rf ${DIR} 2> /dev/null
mkdir -p ${DIR}

openssl genrsa -out ${DIR}/${CA}-priv-key.pem 2048
openssl req -x509 -new -nodes -key ${DIR}/${CA}-priv-key.pem -sha256 -days 800 -out ${DIR}/${CA}-cert.pem \
	-subj "/C=UK/ST=England/L=Staffordshire/O=HFLOW/CN=HFLOW CA"
openssl x509 -pubkey -noout -in ${DIR}/${CA}-cert.pem > ${DIR}/${CA}-pub-key.pem

echo "ca pki files in ${DIR}"