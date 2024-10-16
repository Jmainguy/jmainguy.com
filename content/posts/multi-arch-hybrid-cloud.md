---
title: "Multi Arch Hybrid Cloud"
date: 2024-10-15T12:41:54-04:00
draft: false
---

A new friend of mine wanted to pair program on an interesting challenge, something that was relevant to DevOps and would require a Kubernetes cluster. So in order to prepare for our working session I went ahead and built a cluster.

I knew from past experience [k3s](https://k3s.io/) was easy to setup and was able to scale and add more nodes easily as well, so I started there. 

![Easy](/images/k3s.png)

Some requirements for my new cluster were as follows

* inexpensive, I am quite cheap.
* 24/7 uptime, it should'nt go down if I lose power at my house or my son unplugs things.
* https ingress for web apps I want to host on the cluster.
* secure, I do not want to expose the cluster api or internals to strangers on the internet to mess with.
* Support using old laptops / desktops I have lying around the house as servers.


## Current ingess before k8s
At some point, you have to accept a single point of failure, and money is part of the equation for making that decision. I am already pretty happy with my current ingress solution to all my websites, which happens to be HAProxy + LetsEncrypt certs managed via a flat file and a shell script.

```bash
[root@phy01 ~]# crontab -l
30 22 * * * /root/infra/letsencrypt/cron_certs.sh /opt/certbot /root/infra >> /var/log/certs.log
0 23 * * * /root/infra/scripts/backup.sh >> /var/log/backup.log
[root@phy01 ~]# cat /root/infra/letsencrypt/cron_certs.sh
#!/bin/bash
if [ "$#" -ne 2 ]; then
    /usr/bin/echo "You must enter exactly 2 command line arguments"
    /usr/bin/echo "cron_certs.sh CERTBOTDIR INFRADIR"
    /usr/bin/echo "For example: cron_certs.sh /opt/certbot /opt/infra"
    exit 1
fi
# Podman needs $PATH defined, which we get by sourcing .bashrc
source /root/.bashrc > /dev/null
CERTBOTDIR=$1
INFRADIR=$2
space ()
{ # Print spaces
    /usr/bin/echo ""
    /usr/bin/echo ""
}
# Bak up /etc/hosts so we can check the real connection
/usr/bin/bak -f /etc/hosts
/usr/bin/cp ${CERTBOTDIR}/fake_hosts /etc/hosts
for DOMAIN in $(cat ${CERTBOTDIR}/hostnames); do
    /usr/bin/echo ${DOMAIN}
    ${INFRADIR}/letsencrypt/lets_certs.py ${DOMAIN} ${CERTBOTDIR} ${INFRADIR}
    space
done
/usr/bin/unbak /etc/hosts.bak
```

 I manage DNS via my existing Virtual Private Servers from Buyvm / Linode on the CLI using NSD and bind files.

```bash
vm2		    IN	A	    138.2.152.196
rblack		IN	A	    34.44.160.23
spam		IN	A	    144.76.41.204
spam		IN	AAAA	2a01:4f8:191:4298::2
argocd		IN	A	    144.76.41.204
flights		IN	A	    144.76.41.204
longhorn	IN	A	    144.76.41.204
zot		    IN	A	    144.76.41.204
```

144.76.41.204 being the IP of my Hetzner dedicated server / haproxy host / current docker solutions for websites. It has incredible uptime, and I am quite happy with it. So it being my single point of failure was acceptable.

My house is using Google Fiber, and I have'nt been able to get into the router / network management of it because my wife signed up for it and I don't feel like asking her for her google login. So exposing a port on my home network is impossible with these constraints. That in mind, I wanted to connect my home servers to my Hetzner server on a virtual local area network (VLAN), and I went with a Hub-and-spoke model using wireguard to achieve this. 

## Wireguard
On phy01.standouthoust.com, the hub on hetzner server. 

```bash
[root@phy01 ~]# dnf install wireguard-tools
[root@phy01 ~]# wg genkey | tee /etc/wireguard/$HOSTNAME.private.key | wg pubkey > /etc/wireguard/$HOSTNAME.public.key

[root@phy01 ~]# cat /etc/wireguard/wg0.conf
[Interface]
PrivateKey = NOTMYREALKEY
Address = 10.0.0.1/24
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
ListenPort = 51820

[Peer]
# white-amd
# runs k3s cluster
PublicKey = wWkCFIOn7R39JZ+uihGzBeqkAo+6fQtl4TqMcX1cElY=
AllowedIPs = 10.0.0.2/32
PersistentKeepalive = 25

[Peer]
# rpi4
# runs flightaware
PublicKey = mfwmXlrRCNzKHsPoU8zZE/kEni+/7V2HD7OCPEtOGRo=
AllowedIPs = 10.0.0.3/32
PersistentKeepalive = 25


[root@phy01 ~]# systemctl enable --now wg-quick@wg0
```

And then on my debian based servers, see [why I am using debian for servers now](https://jmainguy.com/posts/enterprise-linux/). I do 

```bash
root@black-intel:/etc/wireguard#  apt install wireguard
root@black-intel:/etc/wireguard#  wg genkey | tee /etc/wireguard/$HOSTNAME.private.key | wg pubkey > /etc/wireguard/$HOSTNAME.public.key

root@black-intel:/etc/wireguard# cat phy01.conf 
[Interface]
PrivateKey = AGAINNOTMYREALKEY
Address = 10.0.0.7/24

[Peer]
PublicKey = Ox7513t/BRVud4Jq32WGlqONFvrqQFZLQutwG1tJaig=
AllowedIPs = 10.0.0.0/24
Endpoint = phy01.standouthost.com:51820
PersistentKeepalive = 25

root@black-intel:/etc/wireguard#  systemctl enable --now wg-quick@phy01
```

This creates the 10.0.0.0/24 network, where all traffic for this network goes over wireguard, and the rest of the traffic uses normal routes out to the internet.

Initially I tried getting phy01.standouthost.com to be the master server in the cluster, but `k3s` does a bunch of firewall work under the hood and it was messing with my existing services, so instead I opted to use my arm server from oracle for this. Which I reinstalled as `k3s.soh.re`

## nftables

We got networking solved right? Well, not quite. We got wireguard vlan working right, but k3s does alot of wacky firewall things using iptables under the hood, so we need to configure nftables as well to ensure we arent listening on 0.0.0.0 for any of this traffic.

Shoutout to [HG from AWX](https://groups.google.com/g/awx-project/c/H8LIrjHhd-c) for coming up with a good nftables ruleset I utilized.

```bash
root@black-intel:/etc/wireguard# cat /etc/nftables.conf 
#!/usr/sbin/nft -f

flush ruleset

table inet jmainguy {
	# protocols to allow
	set allowed_protocols {
		type inet_proto
		elements = { icmp, icmpv6 }
	}

	# interfaces to accept any traffic on
	set allowed_interfaces {
		type ifname
		elements = { "phy01" }
	}

	# ips allowed for k3s internally
	set allowed_cluster_ips {
		type ipv4_addr 
        flags interval
		elements = { 10.0.0.1/24, 127.0.0.1/8, 10.42.0.0/16, 10.43.0.0/16 }
	}

	# services to allow
	set allowed_tcp_dports {
		type inet_service
		elements = { ssh }
	}

	# Chain to mark UDP packets
	chain mark_udp {
		type filter hook output priority 0; # or another priority
		policy accept;

		udp dport 8472 mark set 0x0; # Mark packets on port 8472
	}

	chain allow {
		ct state established,related accept

		meta l4proto @allowed_protocols accept
		iifname @allowed_interfaces accept
		tcp dport @allowed_tcp_dports accept
		ip saddr @allowed_cluster_ips accept
        ip daddr @allowed_cluster_ips accept
	}

	# base-chain for traffic to this host
	chain INPUT {
		type filter hook input priority filter + 20
		policy accept;

		jump allow
		reject with icmpx type port-unreachable
	}

        # Forward chain for inter-node traffic
        chain FORWARD {
                type filter hook forward priority filter + 10
                policy accept;

                jump allow
        }
}
```

I went through alot of variations and finally got it working with the one above, I believe I can go back and optimize this ruleset some more, but it works and I am happy for now. 

## k3s

Great, we got networking, firewall, and ingress completed. Next thing we need to do is actually install the cluster.

I am pretty lazy, so on each node I start with the k3s bash script for getting the systemctl file setup. In the future I can optimize and use ansible instead for new nodes.

```bash
# get systemctl file
root@k3s-soh-re:/home/ubuntu# curl -sfL https://get.k3s.io | sh -

# stop the service, we just wanted the file
root@k3s-soh-re:/home/ubuntu# systemctl stop k3s

# remove any trace of cluster it created
root@k3s-soh-re:/home/ubuntu# rm -rf /var/lib/rancher
root@k3s-soh-re:/home/ubuntu# rm -rf /etc/rancher

# Customize cluster for our preference
# Master server
root@k3s-soh-re:/home/ubuntu# cat /etc/systemd/system/k3s.service
[Unit]
Description=Lightweight Kubernetes
Documentation=https://k3s.io

[Install]
WantedBy=multi-user.target

[Service]
Type=notify
EnvironmentFile=-/etc/default/%N
EnvironmentFile=-/etc/sysconfig/%N
EnvironmentFile=-/etc/systemd/system/k3s.service.env
KillMode=process
Delegate=yes
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNOFILE=1048576
LimitNPROC=infinity
LimitCORE=infinity
TasksMax=infinity
TimeoutStartSec=0
Restart=always
RestartSec=5s
ExecStartPre=/bin/sh -xc '! /usr/bin/systemctl is-enabled --quiet nm-cloud-setup.service 2>/dev/null'
ExecStartPre=-/sbin/modprobe br_netfilter
ExecStartPre=-/sbin/modprobe overlay

# This line ensures nftables is loaded before k3s starts
ExecStartPre=/usr/sbin/nft -f /etc/nftables.conf

# Our wireguard IP for this server is 10.0.0.4. we want all agents to talk to us on this IP
ExecStart=/usr/local/bin/k3s \
    server \
	'--disable=traefik' \
	'--advertise-address=10.0.0.4' \
	'--bind-address=10.0.0.4' \
	'--node-ip=10.0.0.4' \
	'--node-external-ip=10.0.0.4' \
	'--flannel-backend=wireguard-native' \
	'--flannel-external-ip'
```

I was having some issue with the default flannel-backend provided by k3s and moving to wireguard-native solved those for me.

On each additional agent (worker node) we add to the cluster, we follow the same steps, and the systemctl file looks mildly different at the end.

The token is retrieved from the master node at `/var/lib/rancher/k3s/server/token`

```bash
ExecStart=/usr/local/bin/k3s \
        agent --server https://10.0.0.4:6443 --token NOTTHEREALTOKEN \
        '--bind-address=10.0.0.7' \
        '--node-ip=10.0.0.7' \
        '--node-external-ip=10.0.0.7'
```

## Success
```bash
jmainguy@fedora:~/Github/jmainguy.com$ kubectl get nodes -o wide
NAME            STATUS   ROLES                  AGE     VERSION        INTERNAL-IP   EXTERNAL-IP   OS-IMAGE                         KERNEL-VERSION      CONTAINER-RUNTIME
black-intel     Ready    <none>                 2d5h    v1.31.1+k3s1   10.0.0.7      10.0.0.7      Debian GNU/Linux 12 (bookworm)   6.1.0-26-amd64      containerd://1.7.21-k3s2
k3s-soh-re      Ready    control-plane,master   4d23h   v1.31.1+k3s1   10.0.0.4      10.0.0.4      Ubuntu 24.04.1 LTS               6.8.0-1013-oracle   containerd://1.7.21-k3s2
lenovo-laptop   Ready    <none>                 4d22h   v1.31.1+k3s1   10.0.0.6      10.0.0.6      Debian GNU/Linux 12 (bookworm)   6.1.0-26-amd64      containerd://1.7.21-k3s2
white-amd       Ready    <none>                 4d6h    v1.31.1+k3s1   10.0.0.2      10.0.0.2      Debian GNU/Linux 12 (bookworm)   6.1.0-25-amd64      containerd://1.7.21-k3s2
```

![Towers](/images/IMG20241013084544.jpg)

![Spahghetti](/images/IMG20241013084552.jpg)

When wanting to upgrade k3s version in the future, we can wget the binary from a github release https://github.com/k3s-io/k3s/releases and place in `/usr/local/bin/`


## Adding new hostnames for ingress on the cluster

I went with istio gateway / virtualservices for ingress on the cluster. When I want to setup a new service I pick a name, add the DNS entry on my nameservers, add the name to my letsencrypt script and run it to get the certificate, add the name to haproxy

```/bin/bash
   ...lots of other acl rules
   acl host_argocd hdr(host) -i argocd.soh.re
   use_backend istio if host_argocd
   ...lots of other backend rules


backend istio
    balance roundrobin

    server k3s.soh.re 10.0.0.4:80 check
    server white-amd 10.0.0.2:80 check
    server lenovo-laptop 10.0.0.6:80 check
```

Visit my [zot](https://zot.soh.re/) server to see the ingress in action, use Continue as guest when asked to login.


## Postmortem

I am quite proud of the setup, and look forward to moving more services from my traditional docker approach to this kubernetes service.

This is the biggest change I have made to my personal infrastructure in four or five years, and I am excited to go all in on kubernetes.

Some tech debt if I find the time to work on it.

* refactor nftables to only include what is needed
* ansiblize node setup
* move haproxy for k8s to the arm server
* move wireguard hub to the arm server
* setup haproxy in tcp loadbalancing mode and cert-manager on the k8s cluster
  