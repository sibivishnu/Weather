FROM buildpack-deps

RUN apt-get install -y ca-certificates

# ADD ./cacheUpdater . sibi
ADD ./files files/.
ADD ./conf conf/.
ADD ./runUpdater.sh .
RUN chmod +x ./runUpdater.sh

ENTRYPOINT ["./runUpdater.sh"]
