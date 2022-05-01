# This file copied from:
# https://github.com/brave/brave-browser/blob/38ed8b8116acc5199461135307ab219f128dd711/Dockerfile
# and was accompanied by a Mozilla Public License Version 2.0:
# https://github.com/brave/brave-browser/blob/38ed8b8116acc5199461135307ab219f128dd711/LICENSE
FROM debian:buster

# Create non-root user
ARG USER=docker
ARG UID=1000
ARG GID=1000
ARG PW=docker

RUN useradd -s /bin/bash -m ${USER} --uid=${UID} && \
    echo "${USER}:${PW}" | chpasswd

# Install dependencies
RUN export DEBIAN_FRONTEND=noninteractive ;\
    apt-get update -q &&\
    apt-get install -y \
        autogen \
        automake \
        build-essential \
        ca-certificates \
        ccache \
        curl \
        fonts-liberation \
        gconf-service \
        git \
        gperf \
        libappindicator1 \
        libasound2 \
        libatk1.0-0 \
        libbz2-dev \
        libc6 \
        libcairo2 \
        libcups2 \
        libdbus-1-3 \
        libexpat1 \
        libffi-dev \
        libfontconfig1 \
        libgbm-dev \
        libgcc1 \
        libgconf-2-4 \
        libgdbm-dev \
        libgdk-pixbuf2.0-0 \
        libglib2.0-0 \
        libgtk-3-0 \
        libncurses5-dev \
        libnspr4 \
        libnss3 \
        libnss3-dev \
        libpango-1.0-0 \
        libpangocairo-1.0-0 \
        libreadline-dev \
        libsqlite3-dev \
        libssl-dev \
        libstdc++6 \
        libtool \
        libx11-6 \
        libx11-xcb1 \
        libxcb1 \
        libxcomposite1 \
        libxcursor1 \
        libxcursor1 \
        libxdamage1 \
        libxext6 \
        libxfixes3 \
        libxi6 \
        libxrandr2 \
        libxrender1 \
        libxss1 \
        libxtst6 \
        lsb-release \
        pkg-config \
        sudo \
        wget \
        xdg-utils \
        xvfb \
        zlib1g-dev

# Install Node 14
# it is *very likely* this is only relevant for our brave-browser friends ...
RUN set -ex ;\
    curl -fsSL https://deb.nodesource.com/setup_14.x | bash - ;\
    export DEBIAN_FRONTEND=noninteractive ;\
    apt-get update -q ;\
    apt-get install -y nodejs python python-pip ;\
    pip install requests

RUN set -ex ;\
    curl -fsSL https://storage.googleapis.com/golang/go1.18.1.linux-amd64.tar.gz \
        | tar -xzf - -C /usr/local ;\
    /usr/local/go/bin/go version

RUN sed -i'' -e '/^.sudo/s/ ALL$/ NOPASSWD: ALL/' /etc/sudoers && \
    usermod -G sudo $USER

WORKDIR /home/docker
USER $USER
RUN git clone --depth=1 https://chromium.googlesource.com/chromium/tools/depot_tools.git ;\
    echo 'PATH=$HOME/depot_tools:/usr/local/go/bin:$PATH' >> .bashrc ;\
    mkdir chromium_src

ENV BEACON_CC_WRAPPER=ccache
# ENV BEACON_CREDIRECT_DEBUG=true
WORKDIR /home/docker/chromium_src
