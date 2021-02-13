
.PHONY: build
build :
	@-rm -f ./bin/hflow 2> /dev/null
	@go build -o ./bin/hflow

.PHONY: test
test :
	@go test -cover -count=1 ./...

.PHONY: bench
bench:
	@go test -run=XXX -bench=. -benchtime=5s -benchmem ./...

.PHONY: install
install: build
	@-sudo rm -f /usr/local/bin/hflow 2> /dev/null
	@sudo cp ./bin/hflow /usr/local/bin/hflow && sudo chmod +x /usr/local/bin/hflow 
	@echo "hflow successfully installed"

.PHONY: release
release: build
	@./scripts/release.sh

URL_PATTERN=""
STATUS_PATTERN=""
CAPTURE_FILE=""
REROUTE_DIRECTIVE=""

.PHONY: start
start : stop build 
	@./bin/hflow -v=1 -u="${URL_PATTERN}" -s="${STATUS_PATTERN}" -f=${CAPTURE_FILE}
	
.PHONY: stop
stop :
	-@pkill hflow

.PHONY: ca-export
ca-export : stop build 
	@-mkdir ./bin 2> /dev/null
	@./bin/hflow -e=e > ./bin/hflow-ca-export.pem

.PHONY: exec
exec : stub start
	@echo "hello"
#__________________________________________________________________________________________________________________________
#
# Stub server related targets
#__________________________________________________________________________________________________________________________

.PHONY: stub-stop
stub-stop: 
	-@pkill stub

.PHONY: stub
stub : stub-stop
	@-rm ./bin/stub
	@go build -o ./bin/stub ./cmd/stub/
	@./bin/stub -port=8081 -tls=4431 -v=2 2> ./bin/stub-01.log &
	@./bin/stub -port=8082 -tls=4432 -v=2 2> ./bin/stub-02.log &
	@echo "ensure '127.0.0.1 stub-server-01' is in your host file"
	@echo "ensure '127.0.0.1 stub-server-02' is in your host file"

temp :
	@curl -i -k -x http://127.0.0.1:4443 https://duckduckgo.com/?q=these+are+not+the+droids+you+are+looking+for&va=b&t=hc&ia=web

.PHONY: example
example :
	@echo "for the examples to work, ensure hflow is running (make start), the stubs are running (make stub) and the hflow root cert is exported to the bin (make ca-export)"
	@echo "" && echo "press enter to continue or ^c to exit...." && read
	@echo "sending http get request to stub-server via hflow...."
	@curl -i -x http://127.0.0.1:8080 http://stub-server-01:8081/echo/?qs=http-get-qs-data
	@echo "sending http post request to stub-server via hflow...."
	@curl -i -XPOST -d "http-body-data" -x http://127.0.0.1:8080 http://stub-server-01:8081/echo/?qs=http-qs-data
	@echo "" && echo "sending https request to stub-server via hflow while ignoring cert validation...."
	@curl -i -k --tlsv1.3 -XPOST -d "https-body-data" -x http://127.0.0.1:4443 https://stub-server-01:4431/echo/?qs=https-qs-data
	@echo "" && echo "sending https request to stub-server via hflow with hflow root ca configured...."
	@curl -i --cacert bin/hflow-ca-export.pem --tlsv1.3 -XPOST -d "https-ca-body-data" -x http://127.0.0.1:4443 https://stub-server-01:4431/echo/?qs=https-ca-qs-data

#__________________________________________________________________________________________________________________________
#
# PKI related targets. If updating pki files for use in hflow, ensure variables in pem.go files are updated with output
#__________________________________________________________________________________________________________________________

.PHONY: ca
ca :
	@./scripts/gen_ca.sh

.PHONY: cert
cert : # ex: `make cert DOMAIN="stub-server"`. ca certs must be present, if not, run `make ca`
	@./scripts/gen_cert.sh ${DOMAIN}