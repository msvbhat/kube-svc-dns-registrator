# kube-svc-dns-registrator

This service creates a Route53 record with pod IPs for each eligible
service in the cluster. This is useful in environments where each pod has an IP
which is directly routable from other nodes in the network.

For example EKS cluster with AWS VPC CNI by default attaches an ENI to each
pod. So the applications in the pod are directly reachable with the pod IP,
provided VPC network ACLs and Security Groups allow this. So by creating a
Route53 record with pod IPs for the service, the service is accessible from
within the AWS VPC network. This would bypass the kube-proxy and kube-dns.

This service is helpful in situations where a service needs to be accessed from
outside the Kubernetes cluster but only from within the AWS VPC networking.
For example a Lambda function might need to access the service which is not
exposed to outside the cluster via ingress. In those situations, this service
can be used to keep the pods accessible using a DNS name.

## NOTE: This project is still a work in progress
