FROM ubuntu:18.04

# update os software
RUN apt-get update
RUN apt-get -y upgrade
RUN apt-get install -y git

# get repo
RUN git clone https://github.com/andrewmathies/aaaaaaaaaaaaaaa.git

# download go
CMD ["aaaaaaaaaaaaaaa/setup.sh"]

# setup gopath
RUN export GOROOT=/usr/local/go
RUN export GOPATH=/home/ubuntu/server
RUN export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# build and start server
CMD ["aaaaaaaaaaaaaaa/run.sh"]

EXPOSE 22 80 3000
