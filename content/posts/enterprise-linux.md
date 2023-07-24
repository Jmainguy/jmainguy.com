---
title: "Enterprise Linux"
date: 2023-07-24T15:53:23-04:00
draft: false
---

My first attempt at using Redhat Enterprise Linux (RHEL) was back in my silver spring appartment soon after getting married. I was working as an accountant at a small business, and started to get my feet wet by studying for the CompTIA A+ certification. I was big into torrenting at the time and wanted to practice on what I assumed a system administer would need. I downloaded a copy of Exchange Server and RHEL5. I knew of Redhat from [Revolution OS](https://www.youtube.com/watch?v=Eluzi70O-P4). I know I never actually tried using those images I torrented, I ended up running game servers of vnc ontop of ubuntu and opening multiple shells. This worked for awhile.

At some point I started using CentOS and the command line, and got a job down in North Carolina at another small business supporting RHEL4 servers in college bookstores. up2date sucked, yum was such a huge improvement. A few years and jobs later I was working at Redhat, my dream job. I was running fedora as my daily driver and havent looked back since, its a fantastic desktop operating system supported by a very smart and hardworking community. During my time there, CentOS was bought by Redhat. I wasn't even aware that could happen, how could a community operating system be sold?

Redhat was a good steward and nothing really changed, they promised to give CentOS better hardware and support so they would release even faster. There was other clones at the time, Scientific Linux was big, and of course Oracle Linux which kills me. Its one thing to follow the spirit of open source and release a free operating system CentOS, I felt it was another to clone and sell the clone (Oracle Linux).

I have for a little over a decade run all my personal infrastructure on CentOS / Fedora, and been quite happy with it. I think Systemd is a bit much, but sysv scripts were indeed bad to write. CentOS 7 was the last good enterprise linux in my opinon. Recently Redhat has taken hostile actions, including purposfully killing CentOS. They moved it to streams, so it is no longer a bug for bug compatible clone, ruining its entire purpose. From this Alma and Rocky Linux sprung up. I have never liked Rocky's name, you dont want a stable server operating system named Rocky. Yes, I know why they named it Rocky, no it doesn't change my mind. They both looked like viable options and who knows which one would win, the same way Scientific Linux tried to take on CentOS. So about 18 months ago I redid all my infrastructure. I moved some servers to RHEL 8, some to Alma 8, some to Rocky 8, and some to Oracle 8 (I was already using Oracles cloud for free, so I can't be too much of a hypocrite).

After a year you have to resubcribe RHEL 8 free license, lame. It isnt clear you need to do this, or how, updates just stop happening. For this reason, RHEL is not a viable personal choice, either pay for it or don't use it.

Oracle Linux is made by Oracle and I disagree with pretty much everything that company does, I am quite happy to use their free cloud resources though. Their cloud works, the servers run, I have no real complaints about them. So I would not recommend Oracle as an operating system, if you are going to pay for RHEL, go ahead and pay for RHEL.

CentOS is dead, I have no idea what streams is and no desire to learn.

Rocky Linux really has [terrible Press Relations](https://rockylinux.org/news/2023-06-22-press-release/). Every single response starts with Gregory Kurtze who either refers to himself as President of The Rocky Enterprise Software Foundation, or as CEO of Rocky Linux support provider CIQ. I personally don't really like the way he wields so much control over a community and do not wish to partake in that community based on their Press Releases. As such, Rocky Linux is not a viable OS for me.

Alma Linux as a result of Redhat limiting access to source code, is [no longer committed to being a one for one bug for bug clone of RHEL](https://almalinux.org/blog/future-of-almalinux/), so I don't personally see the point.

Running Fedora as a Server I am not super keen on, I like old and stable compared to constant updates and changes. I do love it as a Desktop / Laptop operating system however, and I really hope Redhat doesnt ruin it with their financial control. For this reason I have decided not to switch to Fedora for my servers.

Debian wins. I have always been fond of Debian as a free open source operating system. They have always stuck to their guns, and are the upstream of all the very popular buntu OS's. This weekend I nuked all my serves on oracle's cloud and brough back [soh.re](https://soh.re) running on Ubuntu (I am lazy, Oracle doesn't have a debian OS to choose). I am not personally very fond of Ubuntu, but it was as close to Debian as I could get on Oracle's cloud easily. I now plan to move the rest of my infrastructure to Debian when I find the time to do so.

I was a huge Redhat fanboi, from before I worked there, until well after. It really pains me to see them lose their way under IBM. It really pains me to switch from using their clones to which I have dedicated the majority of my career to learning and mastering. All good things come to an end, and like shadowman did, so has Redhat. Goodbye old friend.
