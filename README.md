# Summary
hflow is a simple, yet powerful, command-line, debugging http/s proxy server.

# Features
hflow exposes the following features via an easy to use interactive CLI

* `Traffic Capture`: Capture all traffic or filter by request url and/or response status
* `Edit & Continue`: Break on requests and/or responses that match a specified url pattern and edit them before they are forwarded and/or returned
* `Request Re-routing`: Route requests destined for one host to another
* `TLS Support`: Decrypt both HTTP and HTTPS traffic. Add HFLOW's root CA certificate into your client's certficate store for seamless HTTPS traffic interception
* `Response Decoding`: Automatically decodes gzip and brotli encoded responses
* `Dump Mode`: A full capture file can be generated for a session by specifying a capture file as a command line argument when a hflow session is started
 
# Why HFLOW?
There are two popular options in the `*nix` http debugging proxy arena; `Charles` and `mitmproxy`. The former is neither free nor command-line based and therefore not comparable to hflow. The latter is a fully featured, and excellent, example of FOSS software. hflow doesn't have any functionality that `mitmproxy` doesn't, and also lacks much that `mitmproxy` does provide. 

But hflow still has a niche.

hflow is primarily a tool to quickly capture traffic on container instances and servers. While `mitmproxy` could do that task, it is a very weighty install with a lot of dependencies which may not be available, or accessible, on a container or server. It also has a steep learning curve, relatively to hflow, for even basic features (*this is not a criticism of `mitmproxy`: more features mean more options and more complexity; hflow focuses on the most commonly used features, and as such, does less*). 

hflow is a single binary: drop the binary on any host and just run it. The interactive CLI is fluid and extremely simple: this document provides examples and instructions, but few would need them to get started. 

When you're finished, to uninstall, delete the binary. 

If `mitmproxy` is `vim` with a bunch of plug-ins; hflow is `nano`: it does the basics well, works anywhere and (almost) anyone can figure it out.

# Installation
Download the appropriate binary for your system from [releases](https://github.com/comradequinn/hflow/releases), add execute permissions and then execute it. 

To make hflow available globally via the `hflow` command; copy, or symlink, the downloaded binary into `/usr/local/bin/` or any other suitable directory available on your `PATH` environment variable.

Alternatively, the scripts below will download and install hflow for you; select the one appropriate for your system and execute it in a terminal:

```bash
# linux on amd 64: amd64
sudo rm -f /usr/local/bin/hflow 2> /dev/null; sudo curl -L "https://github.com/comradequinn/hflow/releases/download/v1.0.0/hflow.linux.amd64" -o /usr/local/bin/hflow && sudo chmod +x /usr/local/bin/hflow
```

```bash
# macOS on apple silicon: arm64
sudo rm -f /usr/local/bin/hflow 2> /dev/null; sudo curl -L "https://github.com/comradequinn/hflow/releases/download/v1.0.0/hflow.darwin.arm64" -o /usr/local/bin/hflow && sudo chmod +x /usr/local/bin/hflow
```

```bash
# macOS on intel silicon: amd64
sudo rm -f /usr/local/bin/hflow 2> /dev/null; sudo curl -L "https://github.com/comradequinn/hflow/releases/download/v1.0.0/hflow.darwin.amd64" -o /usr/local/bin/hflow && sudo chmod +x /usr/local/bin/hflow
```

## From Source
To build and install hflow from source, run the below from a terminal on a machine with `Git` and `Go >=1.19` installed

```
git clone https://github.com/comradequinn/hflow.git && cd hflow && make install
```

This will clone the repo, compile hflow and then copy the resulting hflow binary to `/usr/local/bin`. As this location is typically included in the `PATH` environment variable, hflow should become globally available after the install completes. 

Once the install has completed, you may optionally delete the cloned repo.

Uninstalling hflow is simply a matter of deleting the file `/usr/local/bin/hflow`.

# Usage
The following examples use `curl` to execute requests, see the [later section](#configuring-client-proxies) for help on configuring proxies for other clients. 

## Capturing Traffic
Execute the below to start hflow proxying on the default ports: 

```bash
# terminal 1
hflow
```

hflow will report that it is now proxying traffic.

```
hflow is listening for http traffic on port 8080 and https traffic on port 4443.

press the 'm' key and hit enter to display the menu...

/ proxying...
```

Press `m` and hit enter to display the menu. 

```
_________________________________________________________________________________________

hflow menu
_________________________________________________________________________________________

S - display the current proxy settings
C - write captured traffic to the terminal (optional traffic filters can be applied)
B - set a breakpoint to allow request or response editing
R - reroute requests to a different host
A - display information about hflow
X - exit the menu without making any changes
_________________________________________________________________________________________

enter the required option: 
```

Enter `c` to select writing a traffic capture to your terminal and press enter. 

When prompted to apply request and response filters, press enter to indicate `no filter`. Filters should be provided where there is considerable traffic and only a subset of requests or responses need to be viewed. 

Finally, hit enter to start the capture.

```
capture started

/ proxying...
```
Open a second terminal and execute the below:

```bash
# terminal 2

# -x sets hflow as a proxy for this request only. the url is a duck-duck-go query
# -k instructs curl to ignore certificate errors; certificate errors can be addressed by installing the HFLOW Root CA certificate into the client's CA cert store (see #installing-the-hflow-root-ca-certificate)
curl -i -k -x http://127.0.0.1:4443 https://duckduckgo.com/?q=are+these+the+droids+I+am+looking+for&va=b&t=hc&ia=web 
```

Observe the request and response traffic capture records displayed in `terminal 1`.

```
1     >> GET https://duckduckgo.com:443/?q=are+these+the+droids+I+am+looking+for&va=b&t=hc&ia=web
2      < 200 OK 8614 bytes in body (source: GET https://duckduckgo.com:443/?q=are+these+the+droids+I+am+looking+for&va=b&t=hc&ia=web)
```

Press the `m` key and hit enter to display the menu and this time select `v` to view one of the captured traffic records. Enter `2` when prompted for the record number, which will display the response data

```
_________________________________________________________________________________________

capture detail for record 2:
_ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ 

HTTP/1.1 200 OK
Strict-Transport-Security: max-age=31536000
X-Frame-Options: SAMEORIGIN
Server: nginx
Via: hflow
X-Xss-Protection: 1;mode=block
Expect-Ct: max-age=0
Content-Type: text/html; charset=UTF-8

... omitted for brevity
```

Press enter when finished reviewing the traffic data. 

Note that the menu presented now contains a `|` option, to cancel the terminal-based traffic capture and return to silently proxying (*or whatever configuration was previously active*). Also note the `s` menu option, which can be used to display the active configuration at any time.

```
_________________________________________________________________________________________

hflow menu
_________________________________________________________________________________________

S - display the current proxy settings
V - view the full contents of a captured traffic record
| - stop writing captured traffic to the terminal 
B - set a breakpoint to allow request or response editing
R - reroute requests to a different host
A - display information about hflow
X - exit the menu without making any changes
_________________________________________________________________________________________

enter the required option: 
```

Stop hflow with `CTRL+C`.

## Editing Requests and Responses
Start hflow and navigate to the menu by pressing `m` as shown in the previous section. 

Enter `b` to configure a breakpoint. 

When prompted enter `droids` for the request match text and `1` to indicate that the breakpoint only applies to the responses to requests containing `droids`

```
_________________________________________________________________________________________

hflow menu
_________________________________________________________________________________________

S - display the current proxy settings
C - write captured traffic to the terminal (optional traffic filters can be applied)
B - set a breakpoint to allow request or response editing
R - reroute requests to a different host
A - display information about hflow
X - exit the menu without making any changes
_________________________________________________________________________________________

enter the required option: b
break on traffic where the request matches: droids
break on request only (0), response only (1), both (2): 1

breakpoint configured. hit enter to apply....

```

Press enter to apply the configured breakpoint.

Open a second terminal and execute the below:

```bash
# terminal 2
curl -i -k -x http://127.0.0.1:4443 https://duckduckgo.com/?q=are+these+the+droids+I+am+looking+for&va=b&t=hc&ia=web
```

Observe, shortly, in `terminal 1` that a notification has been written indicating that a breakpoint has been hit. 

Press enter to edit the response and note that a text editor has now taken control of the terminal and is displaying the captured response to the `curl` request. 

Select it in its entirety and delete it (*if using the default editor of vim, `select all` is `ggVG DEL`, you can choose a different editor by setting the `EDITOR` environment variable to your editor of choice*). 

Replace the deleted original response with the below text (*when editing capture files, take care to honour the HTTP specification, specifically, ensure that `Content-Length` is accurate and the text has a trailing new line*):

```
HTTP/1.1 200 OK
Content-Length: 45

These are not the droids you are looking for

```

Save the file and exit the editor. Observe shortly in `terminal 2` that `curl` renders the edited response, rather than the original.

In `terminal 1`, return to the hflow menu and note that it now contains a `/` option, to remove the active breakpoint. Also note the previously mentioned `s` menu option, which can be used to display the active configuration at any time.

```
_________________________________________________________________________________________

hflow menu
_________________________________________________________________________________________

S - display the current proxy settings
C - write captured traffic to the terminal (optional traffic filters can be applied)
/ - remove the active breakpoint
R - reroute requests to a different host
A - display information about hflow
X - exit the menu without making any changes
_________________________________________________________________________________________

enter the required option: 
```

Stop hflow with `CTRL+C`.

## Rerouting Requests
This example uses the network utility `netcat`. This is available on macOS and most Linux distributions and is normally named `nc`. A variant is also available via the `nmap` project named `ncat`. The examples below use the `nc` form.

Start hflow and navigate to the menu by pressing `m` as shown in the previous section. Enter `r` to configure request rerouting. When prompted enter `duckduckgo.com` as the host to reroute traffic **from** and press enter. Then enter `localhost:8081`, when prompted, as the host to route to traffic **to**. 

```
_________________________________________________________________________________________

hflow menu
_________________________________________________________________________________________

S - display the current proxy settings
C - write captured traffic to the terminal (optional traffic filters can be applied)
B - set a breakpoint to allow request or response editing
R - reroute requests to a different host
A - display information about hflow
X - exit the menu without making any changes
_________________________________________________________________________________________

enter the required option: r
enter the host to reroute traffic from: duckduckgo.com
enter the host to reroute traffic to: localhost:8081
rerouting configuration ready. press enter to apply....
```

Press enter to apply the rerouting configuration.

Open a second terminal and execute the below:

```bash
# terminal 2
nc -l localhost 8081 # this starts netcat listening on the specified host and port, any traffic sent there will appear in this terminal
```

Open a third terminal and execute the below:

```bash
# terminal 3
curl -i -x http://127.0.0.1:8080 "http://duckduckgo.com/?q=are+these+the+droids+I+am+looking+for&va=b&t=hc&ia=web" # note this is http not https as nc does not support tls
```

Observe in `terminal 2` that `netcat` has recieved the request instead of the servers behind `duckduckgo.com`. Optionally, return a valid HTTP response by typing it into `terminal 2` or by pasting the below:

```
HTTP/1.1 200 OK
Content-Length: 45

These are not the droids you are looking for

```

Observe in `terminal 3` that the response sent from `netcat` is rendered by `curl`. 

If you are observing the log file (*by default `hflow.log`*), it is likely that it will now contain an error. If so, this is simply due to spaces and new lines not being as required by the HTTP spec due to the difficulty of accurately pasting whitespace into `netcat` and then closing the connection: it is of no concern for this example.

In `terminal 1`, return to the hflow menu and note that it now contains a `\` option, to cancel the active rerouting. Also note the previously mentioned `s` menu option, which can be used to display the active configuration at any time.

```
_________________________________________________________________________________________

hflow menu
_________________________________________________________________________________________

S - display the current proxy settings
C - write captured traffic to the terminal (optional traffic filters can be applied)
B - set a breakpoint to allow request or response editing
\ - cancel rerouting requests to a different host
A - display information about hflow
X - exit the menu without making any changes
_________________________________________________________________________________________

enter the required option: 
```

Stop hflow and `netcat` with `CTRL+C`.

# Usage

## Help
Execute the below to output all configuration options

```
hflow -h
``` 

## Configuring Client Proxies
To route traffic to hflow, configure your HTTP client's proxy address values to `127.0.0.1:[port]` specifying `8080` and `4443` as the `[port]` values for HTTP and HTTPS, respectively (*unless you have overridden these default ports when hflow was started, in which case use those ports instead*). 

In macOS and many Linux distributions, the system proxy settings can be changed globally in the `Settings` UI and, similarly, browsers allow the specification of proxies for all traffic they generate. Many languages and tools also support setting proxies via well known environment variables, as shown below:

```sh
export HTTP_PROXY="http://127.0.0.1:8080"
export HTTPS_PROXY="http://127.0.0.1:4443" # note that https proxies are still initially connected to via http
```

The below example uses curl and sets the proxy inline with `-x`, applying it only to the current request:

`curl -i -XPOST -d "http-body-data" -x http://127.0.0.1:8080 http://example.com/api/resource`

Once your proxy settings are configured, use your client to make HTTP requests and note the captured HTTP traffic in your terminal (*or wherever `stdout` is redirected*).

Hit `[CTRL] + C` to stop hflow 

## Dump Mode: Creating a Capture File
To run in `dump mode` specify a capture file when starting hflow by passing the `-f` flag and with a file name. 

By default, all traffic in the proxy session will be written to this file. Optionally, this output can be tuned by specifying further flags to limit traffic to a specific URL pattern or to output binary data and to truncate bodies at a certain number of bytes. 

The below captures all traffic to `duckduckgo.com`, including binary payloads and limits the body output to 200 bytes; the resulting data is written to a file named `hflow.capture`

```bash
hflow -f "hflow.capture" -b -u="duckduckgo.com" -l=200
```

## Installing the HFLOW Root CA Certificate
To avoid HTTP client warnings relating to the safety of connections to secured domains when proxying HTTPS traffic, you may wish to add the HFLOW Root CA Certificate into your HTTP clients trusted CA certificate collection. Note that this is a potential security risk as the HFLOW Root CA Certificate is freely accessible on the internet. As such, this is undertaken at your own risk and it is advised that you untrust the certificate when not using hflow.

If you wish to proceed, the certficate can be exported in PEM format using the below command. The resulting PEM file can then be loaded directly into your HTTP client's truested CA certificate collection.

```
hflow -e=e > ./hflow-ca.pem
```

# Contributions
Contributions and suggestions are welcome