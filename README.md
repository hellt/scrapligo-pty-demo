Scrapligo is capable of automating CLI based operations over pseudo terminals. A PTY can be exposed by `ssh` or, for example, `kubectl exec`.

The latter option allows to perform bootstrapping automation for networking devices deployed in remote k8s clusters, when direct network access is not possible.

Read more about this in [this blog post](https://netdevops.me/2021/using-scrapligo-to-with-kubectl-or-docker-exec).