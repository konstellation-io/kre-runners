FROM nvidia/cuda:10.2-cudnn7-runtime-ubuntu18.04

# Maintainer of the Dockerfile
LABEL maintainer "Intelygenz - Konstellation Team"

# Fix error installing pipenv:
#   RuntimeError: Click will abort further execution because
#   Python 3 was configured to use ASCII as encoding for the environment.
ENV LC_ALL C.UTF-8
ENV LANG C.UTF-8
ENV TZ=Europe/Madrid
ENV DEBIAN_FRONTEND noninteractive

ARG USER=kre
ARG ID=10001

# Create a non-root user
RUN groupadd ${USER} --gid ${ID} && \
    useradd --home-dir "/${USER}" --create-home --shell "/sbin/nologin" --uid ${ID} --gid ${ID} "${USER}" && \
    mkdir -p /var/log/app 

WORKDIR /tmp

RUN apt-get update && apt-get -yq install --no-install-recommends \
        software-properties-common && \
        add-apt-repository ppa:deadsnakes/ppa -y && \
        apt-get install -yq --no-install-recommends \
            build-essential \
            libsm6 \
            libxext6 \ 
            libxrender-dev \
            python3.9 \
            python3.9-distutils \
            python3-pip \
            python3.9-dev && \
    # Set Python 3.9 as default
    update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.9 1 && \
        update-alternatives --set python3 /usr/bin/python3.9 && \
        update-alternatives --install /usr/bin/python python /usr/bin/python3.9 2 && \
        update-alternatives --set python /usr/bin/python3.9 && \
        update-alternatives --install /usr/bin/pip pip /usr/bin/pip3 3 && \
        update-alternatives --set pip /usr/bin/pip3

COPY Pipfile* ./

RUN python3 -m pip --no-cache-dir install pipenv && \
    pipenv --clear install --system --deploy

RUN rm --recursive --force -- * && \
        apt-get clean && apt-get autoremove -y && \
        rm -rf /var/lib/{apt,dpkg,cache,log}/

WORKDIR /kre

USER ${USER}

COPY ["install_stopwords.py", "src", "./"]

RUN python3 /kre/install_stopwords.py

CMD ["sh","-c","/usr/bin/python3 /kre/src/main.py 2>&1 | tee -a /var/log/app/app.log"]
