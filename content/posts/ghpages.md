---
title: "Github Pages"
date: 2025-04-01T10:38:57-04:00
draft: false
---

If you're looking to deploy a website quickly and easily, GitHub offers a fantastic service called **GitHub Pages**.

[GitHub Pages is a static site hosting service that takes HTML, CSS, and JavaScript files straight from a repository on GitHub, optionally runs the files through a build process, and publishes a website.](https://docs.github.com/en/pages/getting-started-with-github-pages/about-github-pages)

## Why Choose GitHub Pages?

GitHub Pages is incredibly convenient because:
- **No Infrastructure Hassle**: No need to build containers or manage clusters.
- **Reliability**: Backed by a professional team ([status page](https://www.githubstatus.com/)).
- **GitHub Integration**: Seamlessly integrates with the Git workflow you already know—no FTP, SCP, or cPanel required.

You can even automate deployments using **GitHub Actions**. One popular tool for this is [MkDocs](https://www.mkdocs.org/).

---

## Getting Started with MkDocs

**MkDocs** allows you to define your website using Markdown files and YAML configuration. It's:
- [Easy to set up](https://www.mkdocs.org/getting-started/#getting-started-with-mkdocs).
- Perfectly compatible with [GitHub Actions](https://www.mkdocs.org/user-guide/deploying-your-docs/).

### Custom Domains with MkDocs

If you want a custom domain like `docs.jmainguy.com`, make sure to:
1. Add a `CNAME` file with your domain name in the appropriate directory (see [MkDocs documentation](https://www.mkdocs.org/user-guide/deploying-your-docs/#custom-domains)).
2. Avoid overwriting your custom domain during deployments.

---

## Setting Up Custom Domains in GitHub Pages

To configure a custom domain:
1. Go to **GitHub Repo Settings** → **Pages** → Fill in the "Custom Domain" box.
2. Follow [GitHub's guide](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site/managing-a-custom-domain-for-your-github-pages-site#configuring-a-subdomain) to set up a `CNAME` DNS record pointing to `<user>.github.io` (personal account) or `<organization>.github.io` (organization).

![Custom Domain Example](/images/customDomains.png)

---

## Protecting Against Domain Takeovers

### The Risk of Domain Takeovers

Using a `CNAME` without proper verification exposes you to potential domain takeover attacks. Scammers can exploit unverified domains to host malicious content.

![Custom Domain Takeover Example](/images/takeover.png)

### How to Secure Your Domain

To prevent this:
1. Add a **TXT DNS record** for domain verification. 
2. Go to **User/Organization Settings** → **Pages** → **Add a Domain** and follow the instructions in [GitHub's Domain Verification guide](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site/verifying-your-custom-domain-for-github-pages).

---

## What Happens If Your Domain Gets Hijacked?

If your domain is hijacked (e.g., your site starts selling scams), it’s likely because:
- You skipped the TXT verification step.
- MkDocs overwrote your custom domain settings during deployment.

This is what happened to me atleast, and my website looked like this all of a sudden.

![Hijacked Website Example](/images/hijacked.png)


To fix this:
1. Validate your domain by adding the TXT record from your GitHub Pages settings.
2. Update your MkDocs configuration to prevent overwriting.


---

## Final Thoughts

Github Pages is an awesome way to serve websites, ensure you validate the domain to protect from domain hijacking, go deploy a website =).
