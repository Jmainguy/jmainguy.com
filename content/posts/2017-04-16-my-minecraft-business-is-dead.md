---
title: My minecraft business is dead

date: 2017-04-16T23:51:00+00:00
url: /index.php/my-minecraft-business-is-dead/
al2fb_facebook_link_id:
  - 55700847_10103345462668348
al2fb_facebook_link_time:
  - 2017-04-16T23:51:08+00:00
al2fb_facebook_link_picture:
  - custom=https://pbs.twimg.com/profile_images/2789370216/8765b6ba61039a987bdc1b3bc922bdbf_400x400.png
categories:
  - General

---
I started a minecraft business to help pay for hosting my own minecraft server. It blew up and I ended up getting alot of customers really quick. I stopped marketing to get new clients about 2-3 years ago and when I checked today the last of my word of mouth people stopped paying about 3 months ago. So I am cancelling my whmcs license and redirecting the website to the open source code the business generated (which I am also doing a poor job of maintaining).

<a href="https://standouthost.com" target="_blank">standouthost.com</a> now redirects to <a href="https://github.com/standouthost/multicraft" target="_blank">github.com/standouthost/multicraft</a>.

I was able to do this with the following haproxy code

acl host_standouthost hdr(host) -i standouthost.com
  
redirect location https://github.com/standouthost/multicraft if host_standouthost