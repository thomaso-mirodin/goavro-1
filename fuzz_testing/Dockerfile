FROM golang

RUN go get github.com/dvyukov/go-fuzz/go-fuzz

ADD goavro-fuzz.zip /
VOLUME /workdir
CMD go-fuzz -bin /goavro-fuzz.zip -workdir /workdir
