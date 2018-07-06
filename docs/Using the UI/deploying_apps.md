# Deploying Apps from Helm Repositories

Supergiant provides a convenient way to deploy applications to your Kubernetes cluster. In a nutshell, an app is the abstraction of the sum of all resources required to accomplish the purpose of the app. 

These resources may include multiple containers, custom networking setups, storage configurations, and more, all of which can take hours to properly set up. With Supergiant, you can easily deploy applications without bothering about their individual components. The things you should consider are configuration, storage, and network settings, all of which differ depending on the type of application you want to deploy. Supergiant provides application templates via Kubernetes helm repositories that  dramatically simplify their deployment. 

## App and Helm Repositories

A Helm repository is a collection of charts (packages) installable using Helm package manager and containing common configuration for Kubernetes apps such as resource definitions to configure the application's runtime, various dependencies,  communication mechanisms, and network settings. 

The default repository provided by Supergiant is the `/stable` branch of the official [Kubernetes Charts repository](https://github.com/kubernetes/charts) that includes ~160 configured apps. The repository contains well-tested charts that comply with basic technical requirements for running apps in Kubernetes. MongoDB, TensorFlow, NGINX, WordPress, Apache Hadoop are just some examples of popular software available in this repository. 

## Using Helm Repositories to Deploy the App

To deploy the new app from the default or a custom repository, click "**Apps**" in the main navigation menu and then "**Deploy New App**". You'll be transferred to the application list compiled from all available repositories. In the example below, you may see a process of deploying "**Tensorflow-serving**" application retrieved from the stable Kubernetes charts repo.

![](https://s3-ap-southeast-2.amazonaws.com/sg-github-wiki-images/deploying-apps/app-deploy.gif)

Each chart contains an app's editable configuration. You must choose a cluster to which to deploy the app and create the user-friendly name for the deployment. All other parameters and options can be found in the official documentation for the chart (e.g., see configuration options for the  [Tensorflow-serving app](https://github.com/kubernetes/charts/tree/master/stable/tensorflow-serving) we are deploying in the image above). Fill in all necessary fields according to the chart's documentation, and click a "**Submit**" button. If everything is fine, Supergiant will start deploying the app, the time for which varies across applications. As soon as the app is deployed, you'll see its status changed to "**Running**" in your cluster statistics. 

## Adding Custom Repositories

Besides apps from the stable Kubernetes charts repo, Supergiant supports adding private and custom repositories. To add a new repository, select **App & Helm Repositories** under the **Settings** drop-down menu in the upper header. You'll be forwarded to the Helm repositories page, and there you will see the list of current repositories and empty fields to create a new one. 

![](https://s3-ap-southeast-2.amazonaws.com/sg-github-wiki-images/deploying-apps/app-helm-repositories-view.png)

Just name your new repository and enter its URL. After that, Supergiant will add the repository to its register and refresh the apps' list displaying the newly added items. It's as simple as that!

![](https://s3-ap-southeast-2.amazonaws.com/sg-github-wiki-images/deploying-apps/new-repo-helm.gif)

In the example above, we've added an official [Kubernetes Charts incubator repository](https://kubernetes-charts-incubator.storage.googleapis.com/) that includes apps that have not yet met technical requirements to be listed in the `/stable` repo. You might notice that adding a new repository happens almost instantly. 

## What's Next?

- Learn how to use Kubernetes services
- Learn how to attach persistent volumes
- Learn how to work with load balancers.
- Learn how to monitor your cluster from the Supergiant Dashboard
