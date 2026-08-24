---
title: How to setup travis-ci with goreleaser for automated binary golang releases
author: Jonathan Mainguy

date: 2019-10-28T13:15:51+00:00
url: /index.php/how-to-setup-travis-ci-with-goreleaser-for-automated-binary-golang-releases/
categories:
  - General

---
I have been enjoying writing golang applications for a few years now, and I naively assumed other people who wanted to use them could easily build the apps with instructions I provided easy enough.

I discovered that the differences between operating systems is enough that compiling the applications can be burdensome for others. My buddy [@slm][1] (hes kind of a big deal) mentioned I should check out goreleaser and I am glad I did.

I was able to setup a release pretty quickly with it and got [this pretty page][2] with compiled builds for mac, windows, and linux (including rpm and deb). Here is how I did it.

1. Sign up for [travis-ci][3] (its free for opensource projects, choose the right license for your project)

2. Link and authorize it for your projects (travis will walk you through this)

3. Add an API token for GoReleaser so it can update your repository with binary builds and release notes.

4. Add `.travis.yml` and [.goreleaser.yml][5], then configure them for your project. You can check out my examples in the [k8sCapcity project](https://github.com/Jmainguy/k8sCapcity).

5. Cut a release:

```bash
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```

Watch travis spin on your new release, and then see the updated contents. Thank you [goreleaser][6] and [travis-ci][7] for your great projects.

 [1]: https://unix.stackexchange.com/users/7453/slm
 [2]: https://github.com/Jmainguy/k8sCapcity/releases
 [3]: https://travis-ci.com/
 [5]: https://github.com/Jmainguy/k8sCapcity/blob/master/.goreleaser.yml
 [6]: https://goreleaser.com/
 [7]: http://travis-ci.com

> **Editor's note, August 24, 2026:** In my opinion, GitHub Actions eliminated most reasons to use Travis CI, CircleCI, or other hosted CI runners. Even with GitHub's recent hosting issues, I really love GitHub Actions and the ability to self-host runners.
>
> I maintain reusable GitHub Actions workflows for [Go projects](https://github.com/Jmainguy/golang-workflows), [Docker images](https://github.com/Jmainguy/docker-workflows), and [Helm charts](https://github.com/Jmainguy/helm-workflows).
