FROM buildpack-deps AS docker_Sandbox

RUN apt-get install -y ca-certificates

ADD ./cacheUpdater .
ADD ./files files/.
ADD ./conf conf/.
ADD ./runUpdater.sh .
ADD ./credential_file.json .
RUN chmod +x ./runUpdater.sh

ENV GOOGLE_APPLICATION_CREDENTIALS=./credential_file.json
ENV FIREBASE_APPLICATION_CREDENTIALS=./credential_file.json
ENV HTTP_PORT=443
ENV MAX_QUEUE=1024
ENV MAX_WORKER=16
ENV SCP_SERVER_USER=
ENV SCP_SERVER_HOST=
ENV SCP_SERVER_RSA=
ENV DEVICE_REMOTE_FILE_PATH=sandbox.import
ENV DEVICE_LOCAL_TARGET_FOLDER=./files/
ENV PROJECT_ID=lax-gateway
ENV SUBSCRIPTION_NAME=SandboxGeoSubscription
ENV TOPIC_NAME=Geo

EXPOSE 5001

CMD []

ENTRYPOINT ["./runUpdater.sh"]
