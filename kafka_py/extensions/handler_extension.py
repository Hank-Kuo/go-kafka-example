from kafka_py.api.healthz.routers import router as healthz_router


def register_handler(app):
    app.include_router(healthz_router, tags=["healthz"])
    