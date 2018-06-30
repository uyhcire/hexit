FROM golang:1.10

ENV TF_TYPE cpu
ENV TARGET_DIRECTORY /usr/local
RUN curl -L "https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-${TF_TYPE}-$(go env GOOS)-x86_64-1.8.0.tar.gz" | tar -C $TARGET_DIRECTORY -xz
RUN ldconfig

RUN go get github.com/tensorflow/tensorflow/tensorflow/go
RUN go test github.com/tensorflow/tensorflow/tensorflow/go

RUN apt-get update && apt-get install -y python3-pip python3-dev python-setuptools
RUN easy_install -U pip
RUN pip3 install --upgrade tensorflow
RUN pip3 install ipython
