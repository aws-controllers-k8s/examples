# ACK Example Resources

This directory contains examples of individual ACK resources, separated by their
respective service controllers.

The examples provided in this repository do not cover an exhaustive list of
options for each controller. To see all `Spec` and `Status` fields, see our [API
Reference][api-ref]. 

For each example:
- Download or copy the resource definition into a local version on your system
- Modify any property identified by a dollar sign and all caps (eg.
  `$ROLE_NAME`) to meet the description of your intended resource
- Apply the resource to a cluster where the respective ACK controller has been
  installed (using `kubectl apply -f <file_name>`)

[api-ref]: https://aws-controllers-k8s.github.io/community/reference/