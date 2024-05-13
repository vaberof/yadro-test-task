FROM ubuntu:20.04

RUN mkdir -p /opt/app

COPY cmd/yadro-test-task/build /opt/app

COPY examples /opt/app

WORKDIR /opt/app

ENTRYPOINT [ "./main"]