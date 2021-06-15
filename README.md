# mapping

## Description

## Syntax

~~~ txt
mapping {
 
}
~~~

## Examples

Start a server on the default port and load the *whoami* plugin.

~~~ corefile
example.org {
    mapping {
      uri redis://127.0.0.1:6379/1
      prefix _dns_
      suffix
      connect_timeout 10000
      read_timeout 30000
    }
}
~~~

When queried for "dev2.merch8dns.com.bo-api.importer A", CoreDNS will respond with:

~~~ txt
;; QUESTION SECTION:
;dev2.merch8dns.com.bo-api.importer. IN A

;; ANSWER SECTION:
dev2.merch8dns.com.bo-api.importer. 5 IN A      10.245.79.106
~~~

coredns mapping plugin

