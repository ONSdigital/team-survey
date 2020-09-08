# Setting up a new infrastructure for the lambda
​
In the terminal run following commands one by one. You need to do it every time you start a new session.
​
```
​
alias fly="HTTPS_PROXY=localhost:8118 fly"
​
INSTANCE=`gcloud compute instances list --project ons-ci --zones europe-west2-a --format="value(name)" --limit 1 --filter name:gke-`
​
gcloud compute start-iap-tunnel $INSTANCE 32540 --local-host-port=localhost:8118 --project ons-ci --zone europe-west2-a
​
```
​
When the following message appear: "Testing if tunnel connection works. Listening on port [8118]." leave this terminal running and open a new one.
​
In the new terminal run deploy.sh in the following directory: ci/ci-bootstrap/bin (see readme file in ci-bootstrap directory). You need to do it just once.
​
Next run the following script to deploy the pipeline in concourse:
​
```
​
ci/bin/fly_pipeline.sh
​
```
​
To destroy the pipeline and all infrastructure run this command:
​
```
​
ci/bin/destroy.sh
​
```