# Install ACK using Terraform

This directory contains bash scripts and Terraform files for creating a new EKS
cluster, with a managed node group, and installing an ACK controller.

## Prerequisites

- AWS IAM Permissions for creating and attaching IAM Roles
- Installation of required tools:

  - [AWS CLI](https://aws.amazon.com/cli/)
  - [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
  - [Helm](https://helm.sh/docs/intro/install/)
  - [Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli#install-terraform)
  - [eksctl](https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html)

## Terraform Modules

The Terraform modules in this directory use [Amazon EKS Blueprints for
Terraform](https://aws-ia.github.io/terraform-aws-eks-blueprints/main/)

By default, the Terraform modules create the following resources in the
`eu-west-1` region:

- A VPC with 6 subnets (3 Private, 3 Public)
- An EKS Cluster with Kubernetes version set to 1.22
- An EKS Managed Node group

> To modify the default values, edit any files mark with `*.tf`

### Creating an EKS Cluster

Run the following command to create the resources:

```shell
terraform init
terraform plan
terraform apply --auto-approve
```

> PS:
>
> - These resources are not Free Tier eligible.
> - You need to configure AWS Authentication for Terraform with either
>   [Environment
>   Variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html#envvars-set)
>   or AWS CLI [named
>   profiles](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html#cli-configure-profiles-create).

You can connect to your cluster using this command:

```bash
aws eks --region <region> update-kubeconfig --name <cluster_name>
```

> You need to change `region` and `cluster_name` parameters.

### Installing an ACK Service Controller

When you want to install a Service Controller and configure IAM Permissions you
can run `./ack_controller_install.sh <service_name>` and change the
*service_name* accordingly.

The [script](./ack_controller_install.sh) has two functions called install and
permissions.

- Install function downloads the required Helm Chart from the official AWS
  Registry installs it to the Kubernetes cluster.

- Permissions function creates OIDC identity provider for the Kubernetes cluster
  and creates IAM Roles for for Service Accounts of the Service Controllers.

### Cleanup

When you want to delete all the resources created in this repository, you can
run `./cleanup.sh <service_name>` script in the root directory of this
repository and change the *service_name* accordingly.

The [script](./cleanup.sh) has one function and does the following:

- Uninstalls the Helm Chart for Service Controller
- Deletes the CRDs created for Service Controller
- Deletes the OIDC Provider of EKS Cluster
- Deletes the EKS Cluster created with Terraform
