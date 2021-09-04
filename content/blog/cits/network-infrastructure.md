---
title: "Key Networking Concepts and their History"
date: 2021-09-04T13:43:00+11:00
slug: "networking-concepts-history"
description: ""
keywords: ["networking"]
tags: ["networking"]
math: false
draft: true
toc: true
---

Many fundamentals of networking can be confusing to understand when read outside of the context they were defined in.
Several important motivations in understanding why certain protocols or ideas were implemented as-is, are omitted when discussed in isolation.

So, I've collated some key moments in history, and their corresponding concepts, terms, or ideas of that particular moment in time.
Hopefully, when framed in this historic perspective, they become easier to understand and grasp.

## 1969 Croker's Memo: RFCs

{{< tldr "RFCs are memos which provide technical discourse on networking ideas." >}}

Request For Comments (RFCs) are a series of memos, originally intended to be a fast means of distributing ideas with other network researchers.

The first RFC was by Steve Crocker of UCLA, titled "Host Software"[^ha1].

Initially it was intended that several RFC memos would be shared to foster discourse, which when collated and aggregated could then be drafted into some official specification or standard.

Nowadays however, RFCs themselves *are the* official standards and specifications.
Often, you will see some RFC mentioned during discussion of anything pertaining to networking.

## Foundations

### 1969 The ARPA Project: Networks

{{< tldr "Networks are a collection of devices (hosts) that communicate with one another via an established network protocol." >}}

A network is a collection of computers (called "hosts"). The first network, often considered as the conception of the *internet*, is the DARPA funded ARPANET, consisting of multiple connected, high-bandwidth (56 Kbps) links between government, academic, and industrial laboratories[^da1].

ARPANET is the seminal demonstration of host-to-host communication, first achieved with the so-called "1822 protocol", then moving to the "Network Control Protocol (NCP)" in 1970 and eventually to TCP/IP in 1983.

### 1973 Kahn and Cerf on Internetting: Inter-Networks

{{< tldr "An Internetwork is a network of networks; the basis of what the internet is today" >}}

The early concepts of "*the internet*" intended to have multiple individual networks, such as ARPANET, connected and communicating with one another through some "meta-level inter-networking architecture", **irrespective of each individual network's topology and design**.

That is, to form a network of networks - an **internetwork**.

In 1973, Kahn (at DARPA) and Cerf (at Stanford) begun developing a communications protocol to enable this internetworking (then called "internetting"), forming what would later be known as TCP/IP, the foundations of the internet we see and use today.

The DoD would declare TCP/IP to be the official standard for internetting in 1982, and ARPANET would completely migrate to it in 1983.

## A Phonebook for the Internet

### 1981 RFC 791: The Internet Protocol (IP)

{{< tldr "An IPv4 Address identifies a single host and consists of two parts, the network ID to identify which network in the internetwork, and a host identifier to identify the host on that particular network." >}}

The internet protocol (IP) enables inter-networking by detailing a logical addressing and packet delivery system.

Key to the Internet Protocol is the detailing of IP addresses, most notable the IPv4 addresses that has been the dominant addressing scheme used for virtually all of the internet's history to-date.

IPv4 addresses are 32-bit addresses (i.e. thirty-two "1"s and "0"s) that help identify the network a hosts belong to **and the host itself**. These are often seen in their corresponding decimal format, such as `192.168.0.1`, each number separated by a dot representing 8-bits.

### 1983 RFC 882/883: Domain Name System (DNS)

Originally, each computer would resolve memorisable hostnames (e.g. `localhost`) to actual IP Addresses using a simple `HOSTS.TXT` file (see my [past post](https://jamphan.dev/blog/cits/win10-hosts-file/)) that was manually maintained by Stanford Research Institute[^in2].

However, with the migration to TCP in 1983 and consequently the rapid explosion of the internet, issues such as name conflicts arose with this file combined with the high effort to maintain it eventually motived a group including Jon Postel, Paul Mockapetris and Craig Partigethe to develop a new system - the **Domain Name System** in RFC [882](https://datatracker.ietf.org/doc/html/rfc882)/[883](https://datatracker.ietf.org/doc/html/rfc883)[^ha1].

The DNS we hear of today is often thought of as *"the phonebook of the internet"*. It decentralises the address resolution capability across multiple domains (e.g. `.com`, and `.org`) and does so in a hierarchical fashion. DNS is a complex system - you can read more on how it works at [AWS, What is DNS?](https://aws.amazon.com/route53/what-is-dns/)[^aw1].

## Identifying Networks

Once we determined the IPv4 address of a host, we have to locate it with two pieces of information from the address - the network it is a part of (its network number), and its address within the network (the "rest" field -- as in, "the rest of the address").

### Classful Networks
{{< tldr "Classful networking segments the IPv4 address into different classes which each have different lengths of the address reserved for the network identification number." >}}

Originally, the first 8-bits of the IPv4 address were to be reserved for the network number, which allowed 2^8 = 256 independent networks. It seemed sufficient when there weren't many rivals to ARPANET (the 10th network), but it soon become apparent more networks were going to appear.

The first attempt to tackling this issue was **"classful networks"**, devised as part of RFC 791, which expanded the 8-bit network number using classes. The first few bits of an IPv4 address would indicate it's "class" (A, B, or C) which would then inform how many of the 32 bits were then used for the network number. Class A addresses used the first 8 bits, Class B the first 16, and Class C the first 24 bits.

Classful networks were not however *scalable*; One of the biggest issues was the lack of a medium-sized network class. Class C networks (with only 8 bits remaining for hosts) only supported 2^8=256 hosts, whilst the Class B networks supported too many hosts (2^16 = 65536) and it was considered wasteful to hand these limited classes out.

### Subnetting


### 1993 RFC 1519: Supernetting and CIDR

In 1993, RFC 1519 (now in [RFC 4632](https://datatracker.ietf.org/doc/html/rfc4632)) introduced **CIDR (Classless Inter-Domain Routing)** to replace the idea of classful networks

## Configuring Devices

### 1993 RFC 1534: BOOTP and DHCP

## Too many hosts!

### 1996 RFC 1918: NAT & Private Networks

{{< tldr "NAT permits the reuse of IPv4 addresses in different, independent private networks. The effect being that a network of hosts could share one public IP address and internally use a private address." >}}

The limited 32-bit address space continued to be a concern, ultimately this was known as the [IP Address Exhaustion problem](https://en.wikipedia.org/wiki/IPv4_address_exhaustion).

During its conception, it was intended that an IPv4 address was to be given to each unique device.

NAT (Network Address Translation) addresses this by mapping multiple private hosts to one publicly exposed IP address.

## More Reading

For the interested, [internetsociety.org](https://www.internetsociety.org/) has a great article on the brief history of the internet[^in1].

For a detailed timeline of events, check out [cyber.harvard.edu](https://cyber.harvard.edu/icann/pressingissues2000/briefingbook/dnshistory.html)[^ha1].

## References

[^aw1]: (Webpage, n.d.) **AWS**. [What is DNS?](https://aws.amazon.com/route53/what-is-dns/).
[^da1]: (Article, n.d.) **DARPA**. [ARPANET](https://www.darpa.mil/attachments/ARPANET_final.pdf)
[^ha1]: (Webpage, 2000) **Harvard**. [Brief History of the Domain Name System](https://cyber.harvard.edu/icann/pressingissues2000/briefingbook/dnshistory.html)
[^in1]: (Webpage, 1997) **B.M. Leiner, V.G. Cerf, et. al**. [Brief History of the Internet](https://www.internetsociety.org/internet/history-internet/brief-history-internet/), *Internet Society*.
[^in2]: (Webpage, 2016) **K. Meynell**. [Final Report on TCP/IP migration in 1983](https://www.internetsociety.org/blog/2016/09/final-report-on-tcpip-migration-in-1983/). *Internet Society*.
[^sd1]: (Webpage, n.d.) **Science Direct**. [Network Infrastructure](https://www.sciencedirect.com/topics/computer-science/network-infrastructure)

