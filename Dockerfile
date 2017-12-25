FROM jaehue/golang-onbuild
MAINTAINER jang.jaehue@eland.co.kr

# install go packages
RUN go get -u github.com/relax-space/go-kit/... && \
    go get -u github.com/relax-space/lemon-wxpay-sdk && \
    go get -u github.com/relax-space/lemon-wxmp-sdk/... && \
    go get -u github.com/relax-space/lemon-alipay-sdk


# add application
ADD . /go/src/lemon-ipay-api
WORKDIR /go/src/lemon-ipay-api
RUN tar xf tmp/wxcert.tar.gz -C /go/src/lemon-ipay-api
RUN go install

EXPOSE 5000

CMD ["lemon-ipay-api"]