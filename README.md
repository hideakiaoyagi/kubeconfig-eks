# kubeconfig-eks

A command-line tool that constructs 'kubeconfig' file from EKS cluster information.


## Introduction

* I have created this tool to study golang and AWS-SDK.
* I have created this tool because the procedure shown in the [AWS official documents](https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html) is too tiresome.
* Perhaps..., this tool is similar to Azure CLI command `az aks get-credentials`.


## Usage

The simplest usage is as follows:

```
$ kubeconfig-eks --name myekscluster --region us-west-2
```

For more information, run `kubeconfig-eks -h`.
