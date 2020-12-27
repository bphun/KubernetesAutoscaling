FROM ubuntu:18.04

RUN useradd -m proxygen && \
    echo "proxygen:proxygen" | chpasswd && adduser proxygen sudo

RUN apt update && \
    apt install -y wget tar git cmake sudo && \
    wget https://github.com/facebook/proxygen/releases/download/v2020.12.21.00/proxygen-v2020.12.21.00.tar.gz && \
    tar -xzf proxygen-v2020.12.21.00.tar.gz 

# USER proxygen

RUN cd proxygen && \
    ./build.sh

CMD ["cd proxygen"]