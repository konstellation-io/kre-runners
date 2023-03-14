FROM konstellation/kre-py:3.0.0
ENV PYTHONPATH=/kre/.local/lib/python3.9/site-packages/
WORKDIR /tmp
USER root
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install git -y
RUN rm --recursive --force -- * && \
        apt-get clean && apt-get autoremove -y && \
        rm -rf /var/lib/{apt,dpkg,cache,log}/
USER kre
COPY ["Pipfile", "Pipfile.lock", "./"]
RUN pipenv --clear install --system --deploy && \
    pip uninstall nats-py -y && \
    pip install git+https://github.com/konstellation-io/nats.py.git@object_store --no-cache
WORKDIR /kre