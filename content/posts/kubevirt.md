---
title: "Kubevirt"
date: 2024-12-02T11:46:56-05:00
draft: false
---

I believe everything that can be done inside a VirtualMachine can be done just as well or better inside a pod on kubernetes. However, if you insist on using a VM, or have legacy code you do not want to spend time on migrating to Containers, there are several options on the market for making those VM's.

## Broadcom VMWare
Around a year ago, [Broadcom purchased VMWare](https://investors.broadcom.com/news-releases/news-release-details/broadcom-completes-acquisition-vmware). VMWare was the defacto way for enterprise customers to manage VM's on premise before the purchase. Soon after the acquisition, [Broadcom changed the pricing model](https://arstechnica.com/information-technology/2024/10/a-year-after-broadcoms-vmware-buy-customers-eye-exit-strategies/) which in my opinion will be the death of VMWare and creates a unique market opportunity to fill the void left.

## OpenStack
Other options would be Openstack, which provides many of the features the cloud providers do. I personally do not see this catching on and believe it has a perception of being over complicated (I do not have enough experience to say if thats true or not, but it is what I hear people say.)

## RHEV
Another option was Redhat Enterprise Virtualization (RHEV), this was the commercial version of [oVirt](https://www.ovirt.org/) which is itself a wrapper ontop of libvirt (for powering the VM's, and other tools for other features). However, Redhat has decided to abandon RHEV in favor of [Red Hat OpenShift Virtualization](https://www.redhat.com/en/technologies/cloud-computing/openshift/virtualization) (which is a commercial version of [kubevirt](https://kubevirt.io/)).

![RHEV is Dead](/images/rhev.png)

## libvirt

[libvirt](https://libvirt.org/) is a fantastic tool which I have used in combination with [qemu](https://www.qemu.org/) in the past to manage VM's on my hetzner dedicated server. I used [virt-manager](https://virt-manager.org/) as a way of getting graphical console access to the VM's, and installing Linux and Windows VM's. Behind the scenes it uses xml based documents for setting things like the ram, cpu, disk, network, and other settings. I was very happy with this tool when I needed to make VM's and would recommend this path for anyone who doesnt already have a kubernetes cluster.

## kubevirt

Red Hat believes [kubevirt is the future for VM's](https://www.redhat.com/en/blog/unleashing-the-power-of-kubevirt). [Google believes in kubevirt as well](https://cloud.google.com/kubernetes-engine/distributed-cloud/bare-metal/docs/vm-runtime/overview)

kubevirt uses libvirt under the hood, but provides a kubernetes native way of interacting with it.

So with that in mind, I decided this weekend to give a shot myself, and I have to say, I am very impressed.

## My experience with kubevirt

[https://github.com/kubevirt/kubevirt?tab=readme-ov-file#to-start-using-kubevirt](https://github.com/kubevirt/kubevirt?tab=readme-ov-file#to-start-using-kubevirt)  Explains how to get started, and following their kubectl apply -f against the URL's provided really did work well. When I searched for the helm charts I found [This discussion](https://github.com/kubevirt/community/pull/224) where the community is hard at work at providing them. I went ahead and [made my own](https://github.com/Standouthost/helm-charts/tree/main/kubevirt) for now as I wait on upstream to adopt a strategy for helm packaging. 

So now that we have installed kubevirt + [CDI](https://kubevirt.io/user-guide/storage/containerized_data_importer/) (for Persistent data), it is time to actually create a virtual machine.

When I installed debian on my desktops, I downloaded a live iso, burned it onto a thumb drive, and then booted into the thumb drive and followed the graphical installer with a keyboard, selecting various choices like the hostname, partition layout, and passwords. I actually kind of enjoy doing this, however its very not feasible to do via a text console. In the past with libvirt I would use vnc to connect and use a graphical interface and follow similar steps, however for kubevirt I chose to use a [cloud image](https://cdimage.debian.org/images/cloud/) with [cloud-init](https://cloudinit.readthedocs.io/en/latest/reference/examples.html) which allows for automation to run on a new install which allows for quick unattended installations.

I ended up with this cloud config, which 
- Added peristent dhcp (because on kubevirt, the mac changes on destroy / create).
- Installed wireguard and generated a private / public key
- Added my public ssh key so I can ssh into the VM as root once its finished

```bash
#cloud-config
#https://github.com/kubevirt/kubevirt/issues/1646
write_files:
- encoding: b64
    content: bmV0d29yazoKICB2ZXJzaW9uOiAyCiAgZXRoZXJuZXRzOgogICAgaWQwOgogICAgICBkaGNwNDogdHJ1ZQogICAgICBtYXRjaDoKICAgICAgICBuYW1lOiBlbnAqCg==
    owner: root:root
    path: /etc/netplan/99-net-fix.yaml
    permissions: '0644'
packages:
- wireguard
runcmd:
- [ sh, -c, 'wg genkey | tee /etc/wireguard/$(hostname).private.key | wg pubkey > /etc/wireguard/$(hostname).public.key' ]
users:
- name: root
    ssh_authorized_keys:
    - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDgNk9LK/BZk0Sga83brcZMebwgGz0hb+c6u8W8wEMp6uPOVZp2n8KntsB3HAGUPcqY2LXRaMZIzpIftgD1NRlVhxJDyO1rrZmedXA5zy0pGBdkAsnqGXI8ehde4wADEIATWDLHAcgU96TeUs4a4gO/cw7GaXvmubtru5YA0AB/YxiIMTKD+0lHqIUKNrSMIQU4SlzYD4+kdoeBsm57+pFoCUHOB4cZpB3vzHW3IVZOmfDGAQOzSwIM+NlOHFWoyr0RulQXELigSiRzQ0ZRh3BFEieWl/+1nGuxrggqB/x+3I+FpeMMT3RA0mA/an76CmU1ENw5DeF4DELNLf9Q1V9AzzrhLLfdZ89McW68UncsRqEi3sTHYMjljVLtKQtE2M21VEueDDIjUCe1cuPag7AI+YAlTGCysIrCojxzZGQd9a/nh4YxT+GIH0rTEP7un+VwMayschd+OaFxjD+1uTmrBCBg2zFlWxhr+G7juvJuF1FbtraaXY+qEj3BT9YchDE= jmainguy@jmainguy.soh.re
```

I tell kubevirt the size of the disk I want for the VM, as well as the debian image to use via the DataVolume kind

```bash
apiVersion: cdi.kubevirt.io/v1beta1
kind: DataVolume
metadata:
  name: debian-dv
spec:
  source:
    http:
      url: "https://cdimage.debian.org/images/cloud/bookworm/20241125-1942/debian-12-generic-amd64-20241125-1942.qcow2"
  pvc:
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: 20Gi
    storageClassName: longhorn
```

Finally, I create a VirtualMachine kind

```bash
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
  name: debian-vm
spec:
  runStrategy: Always
  template:
    metadata:
      labels:
        kubevirt.io/domain: debian-vm
    spec:
      domain:
        devices:
          interfaces:
            - name: default
              masquerade: {}
          disks:
            - name: disk0
              disk:
                bus: virtio
        resources:
          requests:
            cpu: 1
            memory: 2Gi
      networks:
        - name: default
          pod: {} # Stock pod network
      volumes:
        - name: disk0
          dataVolume:
            name: debian-dv
        - name: cloudinitdisk
          cloudInitNoCloud:
            userData: |
              #cloud-config
              #https://github.com/kubevirt/kubevirt/issues/1646
              write_files:
                - encoding: b64
                  content: bmV0d29yazoKICB2ZXJzaW9uOiAyCiAgZXRoZXJuZXRzOgogICAgaWQwOgogICAgICBkaGNwNDogdHJ1ZQogICAgICBtYXRjaDoKICAgICAgICBuYW1lOiBlbnAqCg==
                  owner: root:root
                  path: /etc/netplan/99-net-fix.yaml
                  permissions: '0644'
              packages:
                - wireguard
              runcmd:
                - [ sh, -c, 'wg genkey | tee /etc/wireguard/$(hostname).private.key | wg pubkey > /etc/wireguard/$(hostname).public.key' ]
              users:
                - name: root
                  ssh_authorized_keys:
                    - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDgNk9LK/BZk0Sga83brcZMebwgGz0hb+c6u8W8wEMp6uPOVZp2n8KntsB3HAGUPcqY2LXRaMZIzpIftgD1NRlVhxJDyO1rrZmedXA5zy0pGBdkAsnqGXI8ehde4wADEIATWDLHAcgU96TeUs4a4gO/cw7GaXvmubtru5YA0AB/YxiIMTKD+0lHqIUKNrSMIQU4SlzYD4+kdoeBsm57+pFoCUHOB4cZpB3vzHW3IVZOmfDGAQOzSwIM+NlOHFWoyr0RulQXELigSiRzQ0ZRh3BFEieWl/+1nGuxrggqB/x+3I+FpeMMT3RA0mA/an76CmU1ENw5DeF4DELNLf9Q1V9AzzrhLLfdZ89McW68UncsRqEi3sTHYMjljVLtKQtE2M21VEueDDIjUCe1cuPag7AI+YAlTGCysIrCojxzZGQd9a/nh4YxT+GIH0rTEP7un+VwMayschd+OaFxjD+1uTmrBCBg2zFlWxhr+G7juvJuF1FbtraaXY+qEj3BT9YchDE= jmainguy@jmainguy.soh.re
```

After applying, you will see a pod get spun up, which downloads the qcow image, and provisions the PVC with it.

After that finishes you will see the VM get spun up

```bash
jmainguy@fedora:~/Github/jmainguy.com$ kubectl get vm
kNAME        AGE     STATUS    READY
debian-vm   4h47m   Running   True
jmainguy@fedora:~/Github/jmainguy.com$ kubectl get vmi
NAME        AGE     PHASE     IP          NODENAME      READY
debian-vm   4h47m   Running   10.42.3.9   black-intel   True
```

I can then ssh into that ip as root

```bash
jmainguy@fedora:~/Github/jmainguy.com$ ssh root@10.42.3.9
Linux debian-vm 6.1.0-28-amd64 #1 SMP PREEMPT_DYNAMIC Debian 6.1.119-1 (2024-11-22) x86_64

The programs included with the Debian GNU/Linux system are free software;
the exact distribution terms for each program are described in the
individual files in /usr/share/doc/*/copyright.

Debian GNU/Linux comes with ABSOLUTELY NO WARRANTY, to the extent
permitted by applicable law.
Last login: Mon Dec  2 09:11:18 2024 from 192.168.86.40
root@debian-vm:~# 
```

I was able to create this VM in minutes using the code, and I can make another ten if I wanted to just as quickly (if I have the disks / resources ) to support it.

## Summary

- Cloud-init is cool
- libvirt is cool
- kubevirt is cool

If you want to create VM's, and already have a kubernetes cluster, I highly recommend kubevirt. You only need two files, a VirtualMachine and a DataVolume. 

If you are looking for a commercial offering to replace VMware, I would recommend offerings which are powered by kubevirt.
