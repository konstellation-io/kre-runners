FROM konstellation/kre-py:3.0.0
ENV PYTHONPATH=/kre/.local/lib/python3.9/site-packages/
WORKDIR /tmp
USER root
COPY nats /usr/bin/
USER kre
WORKDIR /kre