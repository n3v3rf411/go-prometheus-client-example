FROM golang:1.13.3-buster

RUN groupadd --gid 1000 app \
  && useradd --uid 1000 --gid app --shell /bin/bash --create-home app

WORKDIR /go/src/app

USER 1000
RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/" \
  && GO111MODULE=on go get github.com/cespare/reflex@e1e64e8a6dd7552a6e8f616474c91d4d089eada6

CMD ["reflex", "-c", "/reflex.conf"]
