FROM ubuntu:20.04

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

WORKDIR /app
RUN mkdir -p /var/log/app
COPY entrypoint.sh entrypoint.sh

RUN chown -R kre:0 /app \
    && chmod -R g+w /app \
    && mkdir /var/log/app -p \
    && chown -R kre:0 /var/log/app \
    && chmod -R g+w /var/log/app

USER kre

CMD ["sh","-c","./entrypoint.sh 2>&1 | tee -a /var/log/app/app.log"]
