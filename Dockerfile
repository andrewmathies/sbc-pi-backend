FROM ubuntu:18.04

# update os software
RUN apt-get update
RUN apt-get -y upgrade
RUN apt-get install -y git
RUN apt-get install net-tools

# get repo
RUN git clone https://github.com/andrewmathies/sbc-pi-backend.git 

# download go
CMD ["sbc-pi-backend/setup.sh"]

# build and start server
CMD ["sbc-pi-backend/run.sh"]

EXPOSE 80 3000
