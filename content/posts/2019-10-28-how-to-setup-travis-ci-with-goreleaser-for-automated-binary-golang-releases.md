---
title: How to setup travis-ci with goreleaser for automated binary golang releases
author: Jonathan Mainguy
type: post
date: 2019-10-28T13:15:51+00:00
url: /index.php/how-to-setup-travis-ci-with-goreleaser-for-automated-binary-golang-releases/
categories:
  - General

---
I have been enjoying writing golang applications for a few years now, and I naively assumed other people who wanted to use them could easily build the apps with instructions I provided easy enough.

I discovered that the differences between operating systems is enough that compiling the applications can be burdensome for others. My buddy [@slm][1] (hes kind of a big deal) mentioned I should check out goreleaser and I am glad I did.

I was able to setup a release pretty quickly with it and got [this pretty page][2] with compiled builds for mac, windows, and linux (including rpm and deb). Here is how I did it.

Step 1. Sign up for [travis-ci][3] (its free for opensource projects, choose the right license for your project)
  
Step 2. Link and authorize it for your projects (travis will walk you through this)
  
Step 3. Add a api token for goreleaser, to allow it to update your repo with binary builds and notes. I followed this guide <https://goreleaser.com/environment/>
  
Step 4. Add the [.travis.yml][4] and [.goreleaser.yml][5] and configure them to do what you want. You can checkout my examples on the <https://github.com/Jmainguy/k8sCapcity> project.
  
Step 5. Cut a release, 

$ git tag -a v0.1.0 -m &#8220;First release&#8221;
  
$ git push origin v0.1.0

Watch travis spin on your new release, and then see the updated contents. Thank you [goreleaser][6] and [travis-ci][7] for your great projects.

 [1]: https://unix.stackexchange.com/users/7453/slm
 [2]: https://github.com/Jmainguy/k8sCapcity/releases
 [3]: https://travis-ci.com/
 [4]: https://github.com/Jmainguy/k8sCapcity/blob/master/.travis.yml
 [5]: https://github.com/Jmainguy/k8sCapcity/blob/master/.goreleaser.yml
 [6]: https://goreleaser.com/
 [7]: http://travis-ci.com