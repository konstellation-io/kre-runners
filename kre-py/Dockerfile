FROM  nvidia/cuda:10.2-devel-ubuntu18.04

# Create kre user.
ENV USER=kre
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# Fix error installing pipenv:
#   RuntimeError: Click will abort further execution because
#   Python 3 was configured to use ASCII as encoding for the environment.
ENV LC_ALL C.UTF-8
ENV LANG C.UTF-8

RUN apt update && \
    apt -y install libsm6 libxext6 libxrender-dev \
    python3.7 python3-pip python3.7-dev && \
    update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.7 1 && \
    update-alternatives --set python3 /usr/bin/python3.7

WORKDIR /kre

COPY Pipfile* ./
RUN python3 -m pip install pipenv && \
    pipenv install --system --deploy

COPY install_stopwords.py ./
RUN python3 /kre/install_stopwords.py

COPY src src

RUN chown -R kre:0 /kre \
    && chmod -R g+w /kre \
    && mkdir /var/log/app -p \
    && chown -R kre:0 /var/log/app \
    && chmod -R g+w /var/log/app

USER kre

CMD ["sh","-c","/usr/bin/python3 /kre/src/main.py 2>&1 | tee -a /var/log/app/app.log"]
