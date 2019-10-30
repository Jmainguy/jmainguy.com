---
title: I wrote https://soh.re from scratch go check it out
type: post
date: 2017-05-01T17:23:44+00:00
url: /index.php/i-wrote-httpssoh-re-from-scratch-go-check-it-out/
al2fb_facebook_link_id:
  - 55700847_10103383813173628
al2fb_facebook_link_time:
  - 2017-05-01T17:23:52+00:00
al2fb_facebook_link_picture:
  - custom=https://pbs.twimg.com/profile_images/2789370216/8765b6ba61039a987bdc1b3bc922bdbf_400x400.png
categories:
  - General

---
I have always wanted a interactive website that looks like a Linux shell. When @ChrisShort tweeted about GoTTY, I knew I had the essential piece I was missing.

<blockquote class="twitter-tweet" data-lang="en">
  <p lang="en" dir="ltr">
    GoTTY – Share Your Linux Terminal (TTY) as a Web Application &#8211; <a href="https://t.co/qU3r49NLvk">https://t.co/qU3r49NLvk</a>
  </p>
  
  <p>
    &mdash; ChriSSHort ??? (@ChrisShort) <a href="https://twitter.com/ChrisShort/status/856520707609591808">April 24, 2017</a>
  </p>
</blockquote>



I wanted each page reload to bring a fresh copy of the site (because people are mean and will wreck the shell they are in on purpose, I know I would). I wanted it to be secure (I do not want people ddosing others, or bringing down my infra), and I wanted to have a few fun things for them to do (like play crawl, or chat on IRC).

Building the site as a container within Docker, and then adding gotty -w /bin/bash was pretty easy. To get it to load a new container on every page load was hard.

To accomplish that, I built soh-router, which intercepts ws:// websocket calls, checks sqlite3 for an available container, removes that entry from sqlite3, forwards the connection to the new container, and then spawns a new container for the next visitor. It also runs a Reaper periodically to check sqlite3 for containers which are dead now. The code for that can be found at <a href="https://github.com/Jmainguy/soh.re/tree/master/router" target="_blank">Router Code</a>.

I used Haproxy to forward all http requests to a Gotty container which does all the handshakes then tells you to upgrade to websocket (at which point the router intercepts the connection), Iptables to limit outbound connectivity, and Docker to limit jail the user and limit fork bombs.

Was a fun project, I still need to add the rpm, and systemd service file to call this complete, but it is complete enough to share now.

Visit <a href="https://soh.re" target="_blank">https://soh.re</a> and have some fun =).