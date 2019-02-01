# Supporting Teams in Kubernetes
## Overview
One of the main scenarios that we envision with Lyra is being able to support multiple teams creating Workflows at the same time (potentially on the same infrastructure!) from Kubernetes. This proposal will offer several options on how we will achieve this.

## Scenario
As an application service team, I want a way to provide my dev teams an artifact that they can instantiate to manage systems external to Kubernetes using the Kubernetes API.

## Lyra Controller
The purpose of Lyra controller is to enable a pull-based request to provision Workflows from Kubernetes. The Lyra controller listens for Workflow objects being created, and then executes the required lifecycle operations in order to reconcile the desired state of the Workflow.

The controller accepts a `Workflow` CRD object spec ([example](https://github.com/lyraproj/lyra/k8s/attach.yaml)) and provides the following operations:
* Create
* Update
* Delete
* Get
* List

The controller should use the same library as the CLI.

## Architecture
There are two options I think that are viable here:

### Option 1: Centralized controller running on Kubernetes
The first option is that we deploy a Lyra controller on the Kubernetes cluster with a command like `$ lyra init` . The controller would run on the cluster, register the Workflow CRD with the cluster, and start receiving Workflow requests from any client. It would then make the required calls to the resource providers and reconcile state.

<p align="center"><img src="docs/lyra-controller-proposal/kubernetes-based-controller.jpg" alt="kubernetes"></p>

**Advantages of this approach:**
* Supports multiple clients creating Workflow objects
* Only one thing is talking to the identity storage
* Can easily handle concurrent requests
* Kubernetes provides single audit context (i.e. all Workflow objects are registered through K8s)

**Disadvantages of this approach:**
* Single point of failure for availability
* Single point of failure for security
* Needs to be deployed on every Kubernetes cluster
* Lots of weird “registration” features to implement (e.g. Lyra must support a separate set of APIs for registering workflows with a central controller)
* Giant sudo server problem — I have one service that talks to literally every service in the universe with a centralized set of credentials.
* Does not handle different levels of credentials (e.g. I can create EC2 instances, but not RDS, and my colleague is the opposite).

### Option 2: Daemonized controller running on workstation
The second option is that the Lyra controller runs as a daemon on the workstation. In this approach, you could run a command like `$ lyra controller` which would register the Workflow CRD with the Kubernetes cluster and start listening to Workflow creation objects. To ensure that multiple users are not clobbering each other, we would need to implement [Issue #57](https://github.com/lyraproj/lyra/issues/57)  and also provide some mechanism for making sure another listening controller doesn’t try to implement the Workflow (e.g. enforce namespaces?).

<p align="center"><img src="docs/lyra-controller-proposal/workstation-based-controller.jpg" alt="workstation"></p>

**Advantages of this approach:**
* Still supports multiple clients creating Workflow objects
* User-supplied authentication context (e.g. only my AWS keys) are used
* Workflows are already located on my local machine or can be easily download them to my local machine (no registration needed)
* No additional single points of failure (other than the workstation you’re already using)
* Kubernetes still provides single audit context (i.e. all Workflow objects are registered through K8s)

**Disadvantages of this approach:**
* Multiple controllers trying to implement the same workflow
* Each workstation must ensure access to desired resources/resource providers
* Consumers now need to know about how to interact with Lyra 

## Key Considerations
### Team Scaling
In either proposal, Lyra should make it easy to go from a single user and operator to multiple producers AND multiple consumers.

### Identity Management
In either proposal, `identity.db` can live either on the cluster or in a cloud based storage like S3. In both proposals, it must be a shared resource.

###  RBAC on Workflows
In either proposal, to ensure only authorized users can create Workflows, I propose that we leverage the Kubernetes API to implement RBAC controls on Workflow objects. This ensures that Lyra is not creating a separate auth and access control layer.

## Credits
* Large amount of content and ideas borrowed from [Helm v3 proposal]([community/008-controller.md at master · https://github.com/helm/community/blob/master/helm-v3/008-controller.md).
* Key considerations taken from [this article](https://medium.com/virtuslab/think-twice-before-using-helm-25fbb18bc822)
