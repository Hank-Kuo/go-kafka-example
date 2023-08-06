#!/bin/sh

set -e

. /venv/bin/activate

exec python -m kafka_py.consumer