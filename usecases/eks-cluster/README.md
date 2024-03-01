# Build an EKS cluster using only ACK resources

This use case demonstrates how to build an EKS cluster using only ACK resources.
It leverages three different controllers to create the EKS cluster:
- `iam-controller` to create the IAM roles and policies
- `ec2-controller` to create the VPC, Subnets, Route Tables, Internet Gateway, Elastic IP, and Security Groups
- `eks-controller` to create the EKS cluster and NodeGroups.