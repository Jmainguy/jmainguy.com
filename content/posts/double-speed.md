---
title: "Double Speed"
date: 2023-11-06T07:55:01-05:00
draft: false
---

Family was sick this weekend so we took it easy. Leah and I have become big Carolina Hurricanes fans and got season tickets this year. However when they are on the road, the only way to watch them is on Bally. Bally is going thru [bankruptcy proceedings](https://thehockeynews.com/nhl/detroit-red-wings/analysis/bally-bankruptcy-the-future-of-regional-sports-networks) which is what I believe has lead to their [app becoming totally broken](https://www.si.com/fannation/bringmethesports/timberwolves/). When we tried to watch games this week, both nights the app refused to work and we had to watch on some shady gray market websites instead. So, Saturday I spent the day on attempting to set up a whole house vpn, so that the Roku TV devices we use would become part of it, hide our location in atlanta, and attempt to trick espn+ into not blocking out Hurricane games for us.


The first thing I did was grab the new ram I ordered from the last weekend and attempt to install it. Cora Beth assisted me and we got all 4 sticks installed. When I went to boot the computer it powered on, but with a blank screen. I figured it had something to do with the bios / trying to see the new configuration of the ram. After unplugging all but one stick, and then going into bios and enabling the OC profile, I was able to add the sticks back one at a time until it booted up with all 64GB of ram. Did I need to 4x the amount of ram I had? Probably not, but it was fun. I learned after ordering the ram, that indeed, the entire ML model has to fit into ram, preferably GPU ram. This means the ram I bought was pointless for AI / ML experiments I had wanted to play with, and my solution there would be to use 7b or smaller models with a Nvidia GPU. Everything else works on my desktop, and the true ML / AI people are using the cloud for their workloads / training. I believe it would be wasteful for me to buy a card just for that.


So with the ram installed I started looking for Wireguard tutorials as I was hearing it was the new coolest way to do VPN. I have been a big fan of OpenVPN for my needs in the past, but I wanted to see what Wireguard was all about. I found [a good Digital Ocean community blog on it](https://www.digitalocean.com/community/tutorials/how-to-set-up-wireguard-on-debian-11). I have always liked DO community tutorials, they curate / pay for these and they are usually of very high quality. I deleted my old OpenVPN config's and hopped onto my VPS to set up a server. I quickly learned by `cat /etc/*release` that my server was running the now ancient CentOS 7. Again, I quite loved CentOS 7, and the CentOS community in general. However, [I had made a decision to stop using RHEL clones, as I blogged about previously](https://jmainguy.com/posts/enterprise-linux/).

My next course of action was to reinstall the OS, I was able to install the Debian 11 quite quickly using [Linode's ](https://www.linode.com/) web based UI. After logging into the new server, I followed the DO tutorial from earlier and had a Wireguard Server, and Client on my Desktop connecting to it. I was blown away, the speed was basically no different than being off the VPN, the VPN connects / disconnects quickly when asked. When rebooting the server my Desktop automatically reconnects to it. I was sold, no more OpenVPN for me, Wireguard is new best friend.


So now that I got Wireguard working, lets add a client to my router so the whole house can be on the VPN for games. We have a [Netgear Nighthawk R7000](https://www.netgear.com/home/wifi/routers/r7000/), with [Tomato](https://freshtomato.org/) installed. Some quick googling revealed Wireguard is not officially supported, but [could be enabled for Arm based devices](https://wiki.freshtomato.org/doku.php/wireguard_on_freshtomato). I updated Tomato to the latest release and attempted to follow the tutorial. However the very first step of `modprobe` failed. Despite my router being arm, and the tutorial claiming it should work, it did not. I was not going to be able to get a whole house vpn working on this. I googled for some newer routers I could pick up at best buy, and they look quite nice, however I could not really justify spending $300 for a new router I did not really need. Foiled once more.


Leah has been complaining about the wireless access in her office at the front of the house, and with my router no being able to do Wireguard I took a good look at why I was even using it anymore. You see, I have ATT Fiber coming into the house, connected to a ATT Modem / Router combo, which I disabled the wireless on and then connected my Netgear router into. I decided to give the ATT Router another shot by enabling its wireless and unplugged the Netgear, with the idea being I can turn it into a dumb AP and place it in Leah's office once I get a long enough ethernet cord. I went backup stairs to check the connection and now I was getting 2x as fast speeds on fast.com 

![fast.com](https://i.imgur.com/OZHy4UB.png)

I had been plugged into the router before, so apparently the Netgear was limiting me somehow. Cool, double speeds never hurt and now I am getting very close to the speeds advertised / I am paying for.


I went ahead and installed Debian over another VM I had with [BuyVM](https://buyvm.net/) and installed Wireguard there. I then noticed I had not had a third nameserver running when I always wanted 3, so I went into my Oracle Cloud Always free instance and spun up a 4cpu 24gig of ram Arm server with Ubuntu, the closest Oracle lets me get to Debian. I set that up with wireguard and updated me DNS. I now have NYC, Atlanta, and Frankfurt as VPN destinations for my VPN needs.


I rounded out the day by updating the packages on all my vm's and hosts, and upgrading NextCloud from 21 to 27. Running the [Nextcloud Security scanner](https://scan.nextcloud.com/) against my instance gave me a very rewarding feeling ![very rewarding feeling](https://i.imgur.com/CzOUAYs.png)



