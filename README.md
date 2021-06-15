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

    }
}
~~~

When queried for "example.org A", CoreDNS will respond with:

~~~ txt
;; QUESTION SECTION:
;example.org.   IN       A

;; ADDITIONAL SECTION:
example.org.            0       IN      A       10.240.0.1
_udp.example.org.       0       IN      SRV     0 0 40212
~~~

coredns mapping plugin

