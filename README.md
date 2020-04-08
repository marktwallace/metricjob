Metricjob is a project that creates a simple K8s deployment that generates 
Prometheus metrics based on the environment variables in k8s/deployment.yaml.

The job will update it Prometheus metrics every 10 seconds.
So far, there are the following env vars that you can set.

RUN_TIME_IN_SECONDS -- job will terminate after this many seconds

The rest are all command seperated lists. For each type of metric (counter, countervec, etc.) the length of the metric name list must be the same as the other parameters for that type. Example:
```
COUNTER_NAME=apple_grape,banana_nut
COUNTER_RATIO=.9,.25
COUNTER_VEC_NAME=aaa_western,aaa_eastern
COUNTER_VEC_LABEL=state,province
COUNTER_VEC_CARDINALITY=9,5
COUNTER_VEC_RATIO=.5,.8
```

* COUNTER_NAME -- list of names of counter metrics
* COUNTER_RATIO -- probability of an increment each 10 seconds
* COUNTER_VEC_NAME -- list of names of counter vectors (labeled metrics)
* COUNTER_VEC_LABEL -- list of label names, only one label per metric for now
* COUNTER_VEC_CARDINALITY -- cardinality of metrics label, just counting numbers (1,2,3...) for now
* COUNTER_VEC_RATIO -- probability of an increment each 10 seconds


### Building

To build and run metricjob locally, you must have golang installed.

https://golang.org

To build it:

```
cd metricjob
go build
```
We use godotenv to run with a local .env file.

https://github.com/joho/godotenv

To install the godotenv bin command:

`go get github.com/joho/godotenv/cmd/godotenv`

To run locally:

`godotenv ./metricjob`

To build an image to deploy to Kubernetes, you must have docker installed.

https://www.docker.com

The Dockerfile builds an image based on alpine. It used a two step approach where a first image includes dev tools and is used to build a much smaller image for K8s deployment. To build:

`docker build -t metricjob .`

To run the docker image locally:

`docker run --env-file .env -p 8745:8745 metricjob`

Because we use godotenv for local testing, and rely on the deployment.yaml for enviroment in the K8s cluster, we need to specify the .env file to run metricjob in a local docker container. 

To view the Prometheus metrics, you can use a browser to GET the following URL:

`http://localhost:8745/metrics`

### Deploying to K8s

To deploy to K8s, the built docker image must be pushed to a registry that is accessible to the cluster. Do this with docker push. Here is an example using a docker internal repo. Change the URLs to work for your case.
```
docker login docker.internal.myrepo.com
docker tag metricjob:latest docker.internal.myrepo.com/markwallace/metricjob:2
docker push docker.internal.myrepo.com/markwallace/metricjob:2
```

Now that the image is uploaded, you also need imagePullSecrets that are referenced by k8s/deployment.yaml in the namespace where it is deployed. (You ony need to do this the first time.)
```
kubectl -m agent create secret docker-registry metricjob_registry \
    --docker-server=docker.internal.myrepo.com \
    --docker-username=DOCKER_USER \
    --docker-password=DOCKER_PASSWORD \
    --docker-email=DOCKER_EMAIL
```
The K8s deployment and service yaml is in the k8s directory. You will need to modify the deployment yaml to set the env vars as you need them, and to chnage the URL of the metricjob image.

To deploy it all to a K8s cluster,
```
cd k8s
kubectl -n agent apply -f .
```