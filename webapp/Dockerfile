FROM buildpack-deps

RUN apt-get install -y ca-certificates

ADD ./webapp .
ADD ./templates templates/.
ADD ./conf conf/.

ENTRYPOINT ["./webapp"]
