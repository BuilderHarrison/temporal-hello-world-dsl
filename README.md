This was a quick demo I created to just see what type of effort that was required to implement a simple hello world type application in temporal that consumed a JSON file with serverless workflow logic.

My knowledge of serverless workflow probably needs improvement, so happy for corrections and or general comments.

Inspiration and code used from https://github.com/temporalio/samples-go/tree/main/dsl

### Helpful resources:

https://github.com/serverlessworkflow/specification/blob/main/specification.md#Operation-State

https://serverlessworkflow.io/editor.html#

https://github.com/temporalio/temporal/issues/2153#issuecomment-964157514

### General comments

At face value, implementing this for most of our workflows may be very time-consuming, as we would need to be able to create an engine that is capable of understanding all the relevant logic and flows that make makeup a workflow. I feel like, you're almost doubling up the effort to create the serverless workflow and then creating all the relevant handlers/workflows/activities/etc in temporal. Therefore, I think time is best spent doing most of our implementation in pure temporal and then for things like client config we could create a relatively simple engine/service to process the config. If there isn't one already, a JSON schema validator should be created and could be used to ensure the engine is kept in line with the most current serverless workflow templates that clients use, that would be extremely beneficial I believe.

Steps to run this sample:
### Step 1: Clone this repo

In another terminal instance, clone this repo and run this application.

```bash
git clone https://github.com/BuilderHarrison/temporal-hello-world-dsl
cd temporal-hello-world-dsl
```

### Step 2: Run the Workflow

```bash
go run start/main.go
```

### Step 3: Run the Worker

In YET ANOTHER terminal instance, run the worker. Notice that this worker hosts both Workflow and Activity functions.

```bash
go run worker/main.go
```

Next:
1) You can run 
```
go run start/main.go -dslConfig=parallel.json
```
