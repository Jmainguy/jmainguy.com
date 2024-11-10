---
title: "The Case for Ten Gig"
date: 2024-11-10T15:59:36-05:00
draft: false
---

My good friend [Jeff Moore](https://www.linkedin.com/in/jeff-moore-k8s/) challenged me to compare my storage speeds with his homelab setup (knowing his would be perform better). So I obliged and learned a few things in the process, which I think others might find interesting

## [kubestr](https://github.com/kastenhq/kubestr)

[kubestr](https://github.com/kastenhq/kubestr) is a benchmarking tool for kubernetes storage. It works in a clever way by creating the PVC / mounting it inside a pod and running a stress test with nicely formatted results at the end. I was impressed and would recommend it.

Jeff's kubestr score

```bash
Running FIO test (default-fio) on StorageClass (synology-nfs-csi) with a PVC of Size (100Gi)
Elapsed time- 1m13.21783719s
FIO test results:x=674 avg=376.724152
  bw(KiB/s): min=5888 max=86355 avg=48251.207031
FIO version - fio-3.36
Global options - ioengine=libaio verify=0 direct=1 gtod_reduce=1
  blocksize=128k filesize=2G iodepth=64 rw=randwrite
JobName: read_iops
  blocksize=4K filesize=2G iodepth=64 rw=randread
read:
  IOPS=536.328308 BW(KiB/s)=2161
  iops: min=293 max=802 avg=540.172424
  bw(KiB/s): min=1175 max=3209 avg=2162.206787

JobName: write_iops
  blocksize=4K filesize=2G iodepth=64 rw=randwrite
write:
  IOPS=525.518616 BW(KiB/s)=2118
  iops: min=350 max=812 avg=528.448303
  bw(KiB/s): min=1400 max=3249 avg=2115.275879

JobName: read_bw
  blocksize=128K filesize=2G iodepth=64 rw=randread
read:
  IOPS=510.109985 BW(KiB/s)=65825
  iops: min=352 max=640 avg=512.275879
  bw(KiB/s): min=45056 max=82011 avg=65621.687500

JobName: write_bw
  blocksize=128k filesize=2G iodepth=64 rw=randwrite
write:
  IOPS=504.872375 BW(KiB/s)=65158
  iops: min=243 max=612 avg=504.931030
  bw(KiB/s): min=31114 max=78435 avg=64676.484375

Disk stats (read/write):
  -  OK
```

My pre-optimized [kubestr](https://github.com/kastenhq/kubestr) score, testing iscsi mounted storage between the black-intel desktop, and lenovo-laptop in my homelab. Both with a single fast ssd disk. I was unable to get it working using the default 100G test, as my storage only had around 90G left to use.

```bash
jmainguy@fedora:~/tmp$ kubestr fio -s longhorn -z 10G
PVC created kubestr-fio-pvc-9xj86

Pod created kubestr-fio-pod-fwf8r
Running FIO test (default-fio) on StorageClass (longhorn) with a PVC of Size (10G)
Elapsed time- 7m1.388394244s
FIO test results:
  
FIO version - fio-3.36
Global options - ioengine=libaio verify=0 direct=1 gtod_reduce=1

JobName: read_iops
  blocksize=4K filesize=2G iodepth=64 rw=randread
read:
  IOPS=87.566765 BW(KiB/s)=365
  iops: min=14 max=194 avg=90.967743
  bw(KiB/s): min=56 max=776 avg=363.870972

JobName: write_iops
  blocksize=4K filesize=2G iodepth=64 rw=randwrite
write:
  IOPS=46.869923 BW(KiB/s)=203
  iops: min=6 max=164 avg=76.777779
  bw(KiB/s): min=24 max=656 avg=307.166656

JobName: read_bw
  blocksize=128K filesize=2G iodepth=64 rw=randread
read:
  IOPS=90.424545 BW(KiB/s)=12074
  iops: min=30 max=170 avg=94.129028
  bw(KiB/s): min=3840 max=21760 avg=12048.516602

JobName: write_bw
  blocksize=128k filesize=2G iodepth=64 rw=randwrite
write:
  IOPS=53.058487 BW(KiB/s)=7283
  iops: min=2 max=88 avg=58.000000
  bw(KiB/s): min=256 max=11264 avg=7424.633301

Disk stats (read/write):
  sdh: ios=3381/1859 merge=7/119 ticks=1893402/2093660 in_queue=3987061, util=99.554688%
  -  OK
```

So what do all those numbers mean? The read_iops for Jeff averaged at 2162.206787 KiB/s, while mine were 363 KiB/s,  so reading data, my setup was about 7x slower than Jeff's (also, this is comparing my iSCSI against his NFS, iSCSI should be faster than NFS in my uneducated opinion. )

## Network Storage depends on Network

I knew he would beat me, but this was embarrassing and showed something was wrong with my setup. I tried my standard [DD test](https://lowendtalk.com/discussion/42/test-the-disk-i-o-of-your-vps) and was getting only 10.2 MB/s

```bash
16384+0 records in
16384+0 records out
1073741824 bytes (1.1 GB, 1.0 GiB) copied, 105.509 s, 10.2 MB/s
```

Around this time, I discovered that the laptop's Ethernet connection only supported a maximum of 100 Mbps. When you divide this by 8 to convert it to megabytes per second, the theoretical maximum speed is 12.5 MB/s under perfect conditions. Since network storage performance depends on both network and storage speeds to determine IOPS, this meant my setup would never exceed 12.5 MB/s. This was far from ideal, especially considering that my SSD, which I had tested, supported speeds up to 470 MB/s. The 47-fold decrease in speed was quite disappointing.

After confirming the laptop had a usb3 port, which offered speeds of up to 5000 Mpbs, I googled a few usb3 to gigabit adapters

```bash
/:  Bus 03.Port 1: Dev 1, Class=root_hub, Driver=xhci_hcd/4p, 5000M
```

I then remembered I had an extra laptop dock lying around that supported usb3 and gigabit speeds and set that up for the laptop. After configuring `/etc/network/interfaces` for the new dock and swapping some cables around I was able to get 10x speeds on my network storage (you know, 100MiB to 1000MiB is a 10x difference)

```bash
Running FIO test (default-fio) on StorageClass (longhorn) with a PVC of Size (10G)
Elapsed time- 59.649235912s
FIO test results:
  
FIO version - fio-3.36
Global options - ioengine=libaio verify=0 direct=1 gtod_reduce=1

JobName: read_iops
  blocksize=4K filesize=2G iodepth=64 rw=randread
read:
  IOPS=683.792419 BW(KiB/s)=2751
  iops: min=444 max=822 avg=696.233337
  bw(KiB/s): min=1776 max=3288 avg=2785.166748

JobName: write_iops
  blocksize=4K filesize=2G iodepth=64 rw=randwrite
write:
  IOPS=346.631134 BW(KiB/s)=1402
  iops: min=214 max=474 avg=351.366669
  bw(KiB/s): min=856 max=1896 avg=1405.733276

JobName: read_bw
  blocksize=128K filesize=2G iodepth=64 rw=randread
read:
  IOPS=709.734192 BW(KiB/s)=91371
  iops: min=454 max=821 avg=724.166687
  bw(KiB/s): min=58112 max=105170 avg=92709.531250

JobName: write_bw
  blocksize=128k filesize=2G iodepth=64 rw=randwrite
write:
  IOPS=355.914551 BW(KiB/s)=46082
  iops: min=232 max=480 avg=364.033325
  bw(KiB/s): min=29696 max=61563 avg=46627.066406

Disk stats (read/write):
  sde: ios=24225/12181 merge=99/328 ticks=1761911/2075121 in_queue=3837032, util=99.597008%
  -  OK
```

Now my speeds were close to being on par with Jeff's, but still 4.7x slower than my native disks if no network was in the mix.

It is interesting that my reads are 2x as fast as my writes, while Jeff's are even. I assume it has something to do with my disks being optimized for read, but I really don't know.

## Why you need 10 gig networking

Jeff had always talked about his fancy setup with 10 gig networking, and at our previous job the other guys were ordering fiber / copper networking cards for the physical servers (which at the time, I thought was overkill, but fun to have). If I want to get the most out of my storage speeds (470 MB/s) I would need a network that supported up to that speed, (3760 Mib/s), while my current network only supports 1000Mib/s.

This is why 2.5 Gig and 10 Gig networking is a thing, I was right in my previous viewpoint that I would never stream enough video or download an iso from Fedora at speeds that would fill up my 1 gig Google fiber network connection. But I was wrong in my assumption that 2.5 Gig and 10 Gig networking in the datacenter had no real purpose. I know understand it exists for storage.

I looked into upgrading to 10 gig for my home setup briefly, and it would require a new switch, a new PCI card for the desktops, and impossible to add to my laptop. I really enjoy having a hodgepodge of old laptops / desktops for my homelab and this requirement to upgrade would basically mean no more laptops. For this reason, I am going to accept the 1 gig networking limitations, as 100 MB/s speeds are still pretty fast, certainly fast enough for me, for now.


## Thanks Jeff

Thank you [Jeff](https://www.linkedin.com/in/jeff-moore-k8s/) for encouraging me to dig into the setup and learn some more about how it all works
