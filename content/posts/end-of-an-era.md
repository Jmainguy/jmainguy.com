---
title: "End of an Era"
date: 2024-10-30T12:35:16-04:00
draft: false
---

For the past few weeks, I have been migrating all my hosted services over to my new [kubernetes cluster](https://jmainguy.com/posts/multi-arch-hybrid-cloud/) running mainly in my house. It was a good amount of work requiring me to write CI/CD pipelines, helm charts, dockerfiles, goreleaser configs, and learn about OCI annotations, and the limits of iscsi on my local network. When I finished all the work, and the services were all running on k8s, I realized I had no use case to keep around my [hetzner](https://www.hetzner.com/) server any longer. Which makes me sad.

This is a testament to how great a provider Hetzner is. Their servers are fairly priced, offer plenty of bandwidth, great support team, and a really fun and intuitive web interface for managing them. I have been using Hetzner for over a decade, back when I got started with minecraft hosting, and they have never let me down. While I do not need my hetzner server in my new setup, I would love to need to use hetzner again in the future when I can justify a large dedicated server again.

![End of an era](/images/cancel-hetzner.png)

So which services were interesting to move?

## soh.re

If you havent visited [soh.re](https://soh.re) go ahead and check it out. It drops you into the terminal of a docker container, that goes away when you close your browser tab. This is one of my favorite services I host.

I had not visited its infrastructure in some time and came to learn that I had made it quite complicated. I had Haproxy going to Haproxy on each VM, that then shifted traffic to either a static container for issuing the websocket upgrade, or the final docker container for serving the content.

After moving it to my kubernetes cluster, I was able to remove the need for both haproxy and the static websocket upgrade container on each VM, instead running the container on kubernetes for the upgrade. (the payload still happens on VM's outside of kubernetes).

I had to use DestinationRules, and serviceEntries for routing traffic to external services outside the cluster, which was fun to learn about. You can view the code [here](https://github.com/Standouthost/helm-charts/tree/main/sohre/templates) if you are interested in seeing how it is done.

## flights.soh.re

[flights.soh.re](https://flights.soh.re) runs off a raspberry pi in my bedroom for serving a website showing planes flying over my house. So instead of redirecting to the pi from haproxy, now I redirect to it from istio on kubernetes. Similar to how soh.re above works, destinationRules and serviceEntries. Its code is [here](https://github.com/Standouthost/helm-charts/tree/main/flights)


## rightwaypropertycare.com

[rightwaypropertycare.com](https://rightwaypropertycare.com) is my brothers power washing business website. It runs on wordpress, which I learned during this migration project only supports MySQL type databases. I figured everyone supported postgres, but I was wrong. I tried a few different MySQL / MariaDB operators for spinning up mysql clusters, but was underwhelmed with each one I tried. They did not work following their quickstart documentation, which drove me to just spinning up a simple mariadb container as part of the rightway deployment. The website also uses a microservice for combatting spam I wrote. You can see its code [here](https://github.com/Jmainguy/wp-spam). I found out during the migration that it has been down for a few months, which made me sad. I feel much better now that I have it on the k8s cluster as it should not go down again, and will be easier to view if it does.

The database for wordpress needs Persistent storage. Which I had setup previously using longhorn. During the migration I saw I was getting sub 1mbps speeds writing to the storage, which I decided was unacceptable. Digging into the issue, I saw writing directly to the disks using my favorite [dd test](https://lowendtalk.com/discussion/42/test-the-disk-i-o-of-your-vps) I was getting 450mbps speeds on the **black-intel**, but 40mbps speeds on the **lenovo laptop** and **white-amd nodes**. I went ahead and purchased a new drive for each to match what I had in the black-intel. I then reinstalled debian / k3s on the new drives and took out the old spinning slow drives.

![Drives](/images/drives.png)

I also saw that I was trying to use iscsi over wireguard to germany which was 130ms away. A very smart storage engineer at RedHat once explained to me one of his cardinal rules that has always stuck with me, don't send storage across a firewall. So I went into my longhorn UI and disabled the german host from being a storage provider.

![Disabled](/images/disabled.png)

I then also determined to never host anything on the german node that requires persistent storage, which for now was done via cordoning that node so it only handles control plane duties for now. In the future I will use taints / annotations to keep PVC's / amd64 requirements away from it.

After making these changes I went from sub 1mbps to above 300mbps speeds on the PVC's.

**Before**

![Slow Transfer](/images/slow-transfer.png)

**After**

![Fast Transfer](/images/fast-transfer.png)


## etherpad.soh.re

[etherpad.soh.re](https://etherpad.soh.re) is a site I use as a scratchpad / way to share christmas list with relatives easy. It can be used for much more, we used it at Redhat for collaboratively taking notes during meetings.

It had a mariadb database on my german server, and I wanted to migrate it to postgres when moving to the new cluster. I decided on postgres because it had some fantastic operators and people always rave about how much they love postgres. I tried a few different operators but ended up using [cloudnative-pg](https://cloudnative-pg.io/) because its in the CNCF, and supported postgres 17 which some of the others did not. It was pretty simple to setup and I found it to be a great operator.

Doing the migration I learned about another great tool called [pgloader](https://github.com/dimitri/pgloader) This connects to the old mariadb database, and the new postgres database, and then transfers and transforms the data. I really enjoyed using it and was quite impressed it just worked.

## The rest

The rest of websites were pretty straightforward static html websites. I was able to write a [helm chart](https://github.com/standouthost/helm-webapp) they could all use to deliver their content in a uniform manner. This chart is pretty opinionated for my setup using istio, but you can get an idea from it.

Full list of websites and their certificates being provided by LetsEncrypt via cert-manager

![All sites](/images/all-hosts.png)


Shout out to my friend [Sam](https://stackoverflow.com/users/33204/slm) for encouraging me to write more of these blog posts, it means a ton to me.
