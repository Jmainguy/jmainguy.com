---
title: sinatrarb-virt-manager
type: post
date: 2015-07-03T03:52:11+00:00
url: /index.php/sinatrarb-virt-manager/
al2fb_facebook_link_id:
  - 55700847_10102124891258238
al2fb_facebook_link_time:
  - 2015-07-03T03:52:18+00:00
al2fb_facebook_link_picture:
  - media=https://jmainguy.com/?al2fb_image=1
categories:
  - General

---
A client asked me to write him a web app, which would let him stop/start/restart his VM&#8217;s he hosts via libvirt. This is the initial implementation to answer his request.

I next need to add logins for this to be useful at all, but its neat how easy sinatrarb makes creating a web app.

<https://github.com/Jmainguy/sinatrarb-virt-manager> is the code
  
<http://phy01.standouthost.com:8080/> is the demo, go nuts.

I am basically just capturing the action requested from the browser, and passing it to bash.