FROM konstellation/kre-py:3.0.0
ENV PYTHONPATH=/kre/.local/lib/python3.9/site-packages/
WORKDIR /tmp
USER root
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install git -y
RUN curl -L -o nats-cli.deb https://github.com/nats-io/natscli/releases/download/v0.0.35/nats-0.0.35-amd64.deb && \
    dpkg -i ./nats-cli.deb && \
    apt-get install -f && \
RUN rm --recursive --force -- * && \
        apt-get clean && apt-get autoremove -y && \
        rm -rf /var/lib/{apt,dpkg,cache,log}/
USER kre
COPY ["Pipfile", "Pipfile.lock", "./"]
WORKDIR /kre