FROM ubuntu:18.04

# Maintainer of the Dockerfile
LABEL maintainer "Intelygenz - Konstellation Team"

# Fix error installing pipenv:
#   RuntimeError: Click will abort further execution because
#   Python 3 was configured to use ASCII as encoding for the environment.
ARG TZ=Europe/Madrid 
ARG USER=kre
ARG ID=1001
ARG BUILD_ENV

ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

# Create a non-root user
RUN groupadd ${USER} --gid ${ID} && \
    useradd --home-dir "/${USER}" --create-home --shell "/sbin/nologin" --uid ${ID} --gid ${ID} "${USER}" && \
    mkdir -p /var/log/app && \
    chown -R ${USER} /var/log/app

WORKDIR /tmp

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get -yq install --no-install-recommends \
        software-properties-common && \
        add-apt-repository ppa:deadsnakes/ppa -y && \
        DEBIAN_FRONTEND=noninteractive  apt-get install -yq --no-install-recommends \
            curl \
            build-essential \
            libgl1-mesa-glx \
            libsm6 \
            libxext6 \ 
            libxrender-dev \
            python3.9 \
            python3.9-dev \
            python3.9-distutils \
            python3-distutils-extra && \
    curl -L -o nats-cli.deb https://github.com/nats-io/natscli/releases/download/v0.0.35/nats-0.0.35-amd64.deb && \
    dpkg -i ./nats-cli.deb && \
    apt-get install -f && \
    curl -o ./get-pip.py https://bootstrap.pypa.io/get-pip.py && \
    python3.9 ./get-pip.py && \
    # Set Python 3.9 as default
    ln -s /usr/bin/python3.9 /usr/local/bin/python3 && \
    ln -s /usr/bin/python3.9 /usr/local/bin/python

RUN rm --recursive --force -- * && \
        apt-get clean && apt-get autoremove -y && \
        rm -rf /var/lib/{apt,dpkg,cache,log}/

COPY Pipfile* ./

RUN python3 -m pip install --upgrade --no-cache-dir  \
        setuptools \
        wheel && \
    python3 -m pip install --upgrade --no-cache-dir \
        pipenv 


USER ${USER}

RUN pipenv install --system --deploy $([ ! -z $BUILD_ENV ] && echo --dev)

WORKDIR /kre

COPY --chown=kre:0 src /kre/src

CMD ["sh","-c","/usr/bin/python3.9 /kre/src/main.py 2>&1 | tee -a /var/log/app/app.log"]
