FROM ubuntu:latest
LABEL authors="varshim"

ENTRYPOINT ["top", "-b"]