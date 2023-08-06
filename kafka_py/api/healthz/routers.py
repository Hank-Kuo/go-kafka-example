from fastapi import APIRouter

router = APIRouter()

@router.get("/healthz", include_in_schema=False)
async def healthz() -> dict[str, str]:
    return {"status": "ok"}
