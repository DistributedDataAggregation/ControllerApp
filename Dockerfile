FROM ubuntu:latest
LABEL authors="karol"

ENTRYPOINT ["top", "-b"]