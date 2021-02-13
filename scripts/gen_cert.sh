#! /bin/bash

# ex: `./gen_cert.sh "stub-server"`. ca certs must be present, if not, run `./gen_ca.sh`

if [ -z "$1" ]
  then
    echo "specifiy a domain for the certificate"
    exit 1
fi

CA="hflow-ca"
CAFILEROOT="./bin/certs/ca/${CA}"
DIR="./bin/certs/eec"
DOMAIN=$1
read -r -d '' EXT << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = ${DOMAIN}
EOF

rm -rf "${DIR:?}/${DOMAIN}" 2> /dev/null
mkdir -p "${DIR}/${DOMAIN}"
openssl genrsa -out "${DIR}/${DOMAIN}/${DOMAIN}".key 2048
openssl req -new -key "${DIR}/${DOMAIN}/${DOMAIN}".key -out "${DIR}/${DOMAIN}/${DOMAIN}".csr \
	-subj "/C=UK/ST=England/L=Staffordshire/O=HFLOW/CN=${DOMAIN}"
echo "$EXT" > "${DIR}/${DOMAIN}/${DOMAIN}".ext
openssl x509 -req -in "${DIR}/${DOMAIN}/${DOMAIN}.csr" -CA "${CAFILEROOT}-cert.pem" -CAkey "${CAFILEROOT}-priv-key.pem" -CAcreateserial \
	-out "${DIR}/${DOMAIN}/${DOMAIN}.crt" -days 825 -sha256 -extfile "${DIR}/${DOMAIN}/${DOMAIN}.ext"
rm -f "${DIR}/${DOMAIN}/${DOMAIN}.ext"

echo "end entity pki files in ${DIR}/${DOMAIN}"