FROM buildpack-deps

RUN apt-get install -y ca-certificates

ADD ./reporting .
ADD ./templates templates/.

ENTRYPOINT ["./reporting"]
