FROM ubuntu:18.04

# update os software
RUN apt-get update
RUN apt-get -y upgrade
RUN apt-get install -y git

# get repo
RUN git clone https://github.com/andrewmathies/aaaaaaaaaaaaaaa.git

# install go
CMD ["./setup.sh"]

# build and start server
CMD ["./run.sh"]

EXPOSE 22 80 3000
