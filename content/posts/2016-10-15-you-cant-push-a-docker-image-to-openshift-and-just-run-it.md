---
title: You can’t push a docker image to Openshift and just run it

date: 2016-10-15T21:03:43+00:00
url: /index.php/you-cant-push-a-docker-image-to-openshift-and-just-run-it/
al2fb_facebook_link_id:
  - 55700847_10102884374832068
al2fb_facebook_link_time:
  - 2016-10-15T21:03:49+00:00
al2fb_facebook_link_picture:
  - custom=https://pbs.twimg.com/profile_images/2789370216/8765b6ba61039a987bdc1b3bc922bdbf_400x400.png
categories:
  - General

---
So I just finished wasting 3 hours of my life trying to get a docker image to run on Redhats PAAS openshift online, and I failed. In short, you can push a docker image to openshift, but it ignores most of the directives in the Dockerfile and rewrites them as it pleases. It then runs the image as any user it wants, and you cant change that, because you are not a open shift admin, on account of using their PAAS instead of running it yourself. 

And since you cannot control the user, it will get a permission denied if you try and execute any scripts or anything inside the container. Making openshift-online pointless for me personally as a PAAS that wont run my apps. Thanks, no thanks.