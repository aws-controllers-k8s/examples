# ACK Performance at Scale

In this usecase, we would deploy 10K SNS Topics, 20K SQS Queues and 20K SNS Subscriptions i.e 50K resources at a time using SQS and SNS ACK Controllers to examine performance of ACK and other metrics.

**Testing environment:**
* EKS cluster
* ACK SNS and AWS SQS Controller
* Prometheus + Grafana Stack

![](ack.drawio.png)

Prometheus is configured to scrape ACK Controller metrics which will help us analyze the performance and other metrics like resource counts, reconciliation time etc.
*The 50K resources would be deployed using a Helm Chart.*

After installing the Helm chart, following were the observations:

* <u>**Analysing ACK Controller's resources behaviour**</u>
The default installation of SQS and SNS ACK Controllers via Helm Chart configures deployment's CPU and Memory resource limits to 100m and 128Mi respectively. On deploying 50K resources, both SQS and SNS Controllers pods got CRASHED with OOM Killed error.


![](resource_limits_breach.png)

On analysing the pods runtime CPU and Memory resources in Grafana, the controller’s runtime CPU requests were more than the default controller’s CPU Limits (100m), whereas the controllers memory were requesting approximately 400m of CPU and the controller’s default memory limit was configured to be 128Mi, but at the runtime the controller were requesting for approximately 1Gi.

We then upgraded ACK Controllers CPU and memory resource limits to 500m and 1Gi.
```
export SERVICE=sqs
export RELEASE_VERSION=$(curl -sL https://api.github.com/repos/aws-controllers-k8s/${SERVICE}-controller/releases/latest | jq -r '.tag_name | ltrimstr("v")')
export ACK_SYSTEM_NAMESPACE=ack-system
export AWS_REGION=eu-west-1

aws ecr-public get-login-password --region us-east-1 | helm registry login --username AWS --password-stdin public.ecr.aws
helm upgrade -i --create-namespace -n $ACK_SYSTEM_NAMESPACE ack-$SERVICE-controller \
  oci://public.ecr.aws/aws-controllers-k8s/$SERVICE-chart --version=$RELEASE_VERSION --set=aws.region=$AWS_REGION --set=metrics.service.create=true --set=resources.limits.memory=1Gi --set=resources.limits.cpu=500m
```

We did helm install again to deploy 50K resources, and now the controllers requests were way within the limits and there were no visible crashes on the controllers pods.

![](upgraded_resource_limits.png)

* <u>**Analysing time taken to deploy 50K resources**</u>
Now, we were able to deploy 50K resources using Helm chart with upgraded ACK Controllers resource limits. It took 2 hours and 15 minutes to deploy and provision 50K resources. Thats quite a huge amount of time. 

![](deployed_resource_time.png)

defaultMaxConcurrentSyncs is an ACK Controller's parameter which can be set in Controller's Helm Chart values file during installation, is the default number of concurrent syncs that the ACK reconciler can perform. By default, the value of defaultMaxConcurrentSyncs is 1. 

We then set defaultMaxConcurrentSyncs to 200 and tried deploying 50K resources at a time to analyse the time it takes now for provisioning.
```
export SERVICE=sqs
export RELEASE_VERSION=$(curl -sL https://api.github.com/repos/aws-controllers-k8s/${SERVICE}-controller/releases/latest | jq -r '.tag_name | ltrimstr("v")')
export ACK_SYSTEM_NAMESPACE=ack-system
export AWS_REGION=eu-west-1

aws ecr-public get-login-password --region us-east-1 | helm registry login --username AWS --password-stdin public.ecr.aws
helm upgrade -i --create-namespace -n $ACK_SYSTEM_NAMESPACE ack-$SERVICE-controller \
  oci://public.ecr.aws/aws-controllers-k8s/$SERVICE-chart --version=$RELEASE_VERSION --set=aws.region=$AWS_REGION --set=metrics.service.create=true set=reconcile.defaultMaxConcurrentSyncs=200 --set=resources.limits.memory=1Gi --set=resources.limits.cpu=500m
```
Guess what, now it just took 52 minutes to deploy 50K resources.

* <u>**Analysing reconcilation time while creating new additional resources**</u>

The final analysis was to test how much time the ACK Controllers would take to deploy  **500 additional resources** including 100 SNS Topics, 200 SQS Queues and 200 SNS Subscriptions on top of already existing 50K resources. 
This would help us examine, on how ACK deals with reconcilation. 

So we installed another Helm chart which would deploy 500 additional resources. 

![](deploy_additional_resources.png)

Guess what, it just took **30 seconds** to deploy 500 additional resources, and hence the reconcilation time is negligible. 