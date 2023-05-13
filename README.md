Project Overview
===================
The weather service wrapper that provides data from accuweather and other sources(?), with result caching to improve response times and reduce costs. 

# Components 
- Weather App: .... 
- Reporting: ????
- cacheUpdater: ????
- common: common libs shared across projects
- test: ....

# Management
Until K8S issues are resolved we are temporarily running and monitoring this service on the Stage3  server under `/mnt/data/code/weather-service` managed by pm2: 
Until permissions/accounts are setup you will need to sudo su as 'keith_brings' to check/start/build etc.

`pm2 list`
`pm2 k8s-fallback-ws logs`

From the repo root folder `pm2 start k8s-fallback-ws.json` may be used to start the service. Restarts are implemented already with exceptional fallback. 

Environment vars are set by `/mnt/data/code/weather-service/.env`,  go lang version is specified in `./.tools-versions` using asdf. 

For reasons a symlink has been generated mapping webapp/templates to /templates any server running this service. 



# Configuration 

## Environment Variables 

### CacheUpdater 
| ENV                          | Description                      | Notes                                      |
|------------------------------|----------------------------------|--------------------------------------------|
| ACCU_API_KEY                 | API key for AccuWeather API      |                                            |
| REDIS_HOST                   | Redis server host                |                                            |
| HTTP_PORT                    | HTTP server port                 |                                            |
| MAX_QUEUE                    | Maximum number of queued workers |                                            |
| MAX_WORKER                   | Maximum number of worker threads |                                            |
| SCP_SERVER_HOST              | SCP server host                  |                                            |
| SCP_SERVER_USER              | SCP server username              |                                            |
| SCP_SERVER_RSA               | SCP server RSA private key file   |                                            |
| DEVICE_REMOTE_FILE_PATH      | Remote file path on device        |                                            |
| DEVICE_LOCAL_TARGET_FOLDER   | Local target folder for download |                                            |
| PROJECT_ID                   | Google Cloud project ID          | Used for pub/sub and attribute pub/sub      |
| SUBSCRIPTION_NAME            | Pub/Sub subscription name        | Used for receiving messages from topic      |
| TOPIC_NAME                   | Pub/Sub topic name               | Used for sending messages to subscribers   |
| ATTRIBUTE_TOPIC_NAME         | Pub/Sub attribute topic name     | Used for sending attribute updates         |



### WebApp
| ENV                        | Description                           | Notes                                      |
|----------------------------|---------------------------------------|--------------------------------------------|
| ENV_ACCU_API_KEY           | ACCU API key                          |                                            |
| ENV_REDIS_HOST             | Redis host                            |                                            |
| FLAG_HTTP_PORT             | HTTP port                             |                                            |
| FLAG_HTTP_HOST             | HTTP host                             |                                            |
| FLAG_HTTP_SCHEME           | HTTP scheme (http or https)           |                                            |
| ENV_FIREBASE_SERVICE_FILE  | Firebase application credentials file | Path to the JSON file for service account. |

### Summary Data (yaml)






# Infra 
Details on K8N config and deployment pipeline. 



# Sandbox Setup

## Linux
0. Deps 
    Install Docker, GoLang and Redis. Start Redis Server. 

    Current production version of GO is 1.8.7: https://medium.com/@patdhlk/how-to-install-go-1-8-on-ubuntu-16-04-710967aa53c9
    

1. Set /etc/environment Variables
    ```
    ACCU_API_KEY="ACCU_KEYHERE"
    REDIS_HOST=localhost:6379
    HTTP_PORT=443
    HTTP_HOST=https://ingv2.lacrossetechnology.com/
    HTTP_SCHEME=https
    FIREBASE_APPLICATION_CREDENTIALS=/mnt/c/Github/lax-alerts/src/ingressor/src/priv/google/pub-sub.json
    ```

2. Compile and run
    ```
    cd webapp
    make compile
    ./tmp/webapp
    ```



    
3. To avoid docker
### start redis
redis-server

### compile and run application (webapp)
export GOPATH=$(pwd)/_vendor
go build -ldflags '-X main.BUILD=wip' -o ./webapp/tmp/webapp ./webapp/*.go
./webapp/tmp/webapp    

### compile and run application (cacheUpdater)
export GOPATH=$(pwd)/_vendor
go build -ldflags '-X main.BUILD=wip' -o ./cacheUpdater/tmp/cacheUpdater ./cacheUpdater/*.go
./cacheUpdater/tmp/cacheUpdater    


### One Liners 
export GOPATH=$(pwd)/_vendor
#### Build and Run CacheUpdate
cd $GOPATH/../ && go build -ldflags '-X main.BUILD=wip' -o ./cacheUpdater/tmp/cacheUpdater ./cacheUpdater/*.go && cd cacheUpdater/tmp/ && ./cacheUpdater
#### Build and Run webapp
cd $GOPATH/../ && go build -ldflags '-X main.BUILD=wip' -o ./webapp/tmp/webapp ./webapp/*.go && ./webapp/tmp/webapp


## Windows
0. Deps 
    Install Docker, GoLang and Redis for Windows. Start Redis Server.      
1. Set Environment Variables (using non 192.168.0.1/127.0.0.1 ip address from `ipconfig`)
    ```
    ACCU_API_KEY="ACCU_KEYHERE"
    REDIS_HOST=192.168.23.21:6379
    ```
2. Build and Run 
    ```
    cd webapp
    make.cmd build
    make.cmd run
    ```
3. Browse to localhost:5000
