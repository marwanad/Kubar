# Kubar

## Overview

Kubar is a quick command line tool to dump and replay the state of a Kubernetes cluster. It is mainly useful in the cases when there is a need to replicate the environment on a different cluster or for migrations.

Note: This targets stateless applications running on Kubernetes clusters and does not support volume snapshotting (you'd have to do that manually yourself).

## Dependencies
Kubar requires an installation of [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

## Build and Installation

Run `make` to build to kubar/kubar directory and `make install` to install to $GOPATH/bin.

## Sample Usage

To export Kubernetes resources as YAML config files:

```
kubar --path=/output-path --mode=export
```

This outputs Kubernetes configuration files for global resources such as namespaces, storage classes and custom resource definitions.

It also outputs other Kubernetes resources such as deployments, secrets, services, etc for every namespace.

**Note: This also exports secrets as configuration files.**


The exported resources follow the below file structure:

```
.
├── kube-system
│   ├── configmaps.yaml
│   ├── daemonsets.yaml
│   ├── deployments.yaml
│   ├── replicasets.yaml
│   ├── role.yaml
│   ├── rolebinding.yaml
│   ├── secret.yaml
│   ├── services.yaml
│   ├── storageclasses.yaml
├── kube-public
│   ├── configmaps.yaml
│   ├── role.yaml
│   ├── rolebinding.yaml
│   └── storageclasses.yaml
└── default
|   ├── services.yaml
|   ├── storageclasses.yaml
└── deployments
|   ├── deployments.yaml
|   ├── secret.yaml
|   ├── services.yaml
|   ├── storageclasses.yaml

```

To apply or replay those configurations, create a new Kubernetes cluster and ensure that the correct cluster role bindings are set to be able to apply the changes. Then simply run:

```
kubar --path=/backup-path --mode=restore
```

This would start by recovering the exported global resources and then applies the configuration files per namespace starting with kube-system and finally going through the custom ones.