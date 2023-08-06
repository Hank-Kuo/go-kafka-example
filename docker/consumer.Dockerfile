FROM python:3.9-slim as base

WORKDIR /app 

ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PIP_NO_CACHE_DIR=off \
    PIP_DISABLE_PIP_VERSION_CHECK=on \
    PIP_DEFAULT_TIMEOUT=100 \
    POETRY_VIRTUALENVS_IN_PROJECT=true \
    POETRY_NO_INTERACTION=1 


FROM base as builder

ENV POETRY_VERSION=1.5.1

RUN pip install "poetry==$POETRY_VERSION"
RUN python -m venv /venv

COPY . ./
RUN poetry export -f requirements.txt | /venv/bin/pip install -r /dev/stdin
RUN poetry build && /venv/bin/pip install dist/*.whl


FROM base as final
COPY --from=builder /venv /venv
COPY . /app
RUN chmod +x ./docker/docker-consumer-entrypoint.sh
CMD ["./docker/docker-consumer-entrypoint.sh"]

EXPOSE 8001