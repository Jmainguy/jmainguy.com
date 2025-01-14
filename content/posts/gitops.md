---
title: "GitOps: Transforming Deployment with Automation and Version Control"
date: 2025-01-14T09:53:35-05:00
draft: false
---

# What is GitOps?

GitOps is a way to manage infrastructure and application deployments using Git as the single source of truth. It focuses on:

* **Declarative configurations:** The desired state of your system (infrastructure or applications) is defined in Git repositories.
* **Automation:** Changes to the Git repository trigger automated updates to the system.
* **Reconciliation loops:** A GitOps operator (e.g., ArgoCD) continuously ensures the system's actual state matches the desired state in Git.

## Key Features of GitOps

* **Version-controlled changes:** Every change to the system is tracked in Git, providing auditability and rollback capabilities.
* **Automation:** Tools like ArgoCD automate applying changes from Git to the actual environment.
* **Reconciliation:** Continuous monitoring ensures the system stays aligned with the desired state in Git.
* **Rollback:** Easily revert to a previous state by rolling back changes in Git.


## My experience before GitOps

During my time at [Red Hat](https://www.redhat.com/), I was responsible for deploying application code to various environments. This process involved several manual steps:

1. The application team documented their [RPMs](https://en.wikipedia.org/wiki/RPM_Package_Manager) in a wiki.
2. I gathered these RPMs and used a tool called **Juicer** to package them in an artifact system.
3. Using an Ansible-like system called **taboot**, I:
    - Removed half the servers from the load balancer.
    - Performed the upgrade on those servers.
    - Paused to verify the upgrade before proceeding with the remaining servers.

While this workflow was effective, it was labor-intensive, requiring highly skilled engineers. We called ourselves **Release Lieutenants**, a title that came with physical badges. The platform team designated a **Captain** who trained and supported us. 

Although this method fostered a strong team culture, it relied heavily on manual intervention and expertise. In hindsight, adopting GitOps could have saved time and resources by automating much of this process, allowing Red Hat to allocate more resources to application development.


## My experience with GitOps

Today, my blog serves as a practical example of GitOps in action, demonstrating how this approach simplifies and automates deployments. Here’s how it works:

1. After previewing a change locally, I merge and push the code to [GitHub](https://github.com/jmainguy/jmainguy.com).
2. A GitHub Webhook triggers an update process on my server, which performs a fresh `git pull`.
3. The [Kubernetes YAML](https://github.com/Standouthost/clusters/blob/main/k3s.soh.re/jmainguy.com.yaml) configuration for the website is synchronized with [ArgoCD](https://argo-cd.readthedocs.io/) and deployed to the cluster.

Even the [Golang code](https://github.com/Jmainguy/autoweb/blob/main/main.go) that automates the git pull process is managed with Git. It’s version-controlled, built into a container, and stored in an artifact registry, ready to deploy.

With this streamlined setup, deploying a new post is as simple as:

```
hugo
git add .
git commit -m 'feat: new posts'
git push
```

I manage all my Kubernetes applications using GitOps because I value seamless, reliable deployments that can be easily rolled back if necessary. Declarative configurations lie at the core of this philosophy, ensuring systems are intuitive and understandable simply by reading the code.
