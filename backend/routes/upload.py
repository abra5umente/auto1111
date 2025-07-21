from fastapi import APIRouter, UploadFile, File, HTTPException
from fastapi.responses import JSONResponse
from pathlib import Path
import logging
import json

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

router = APIRouter()

UPLOAD_DIR = Path(__file__).resolve().parent.parent / "settings"
UPLOAD_DIR.mkdir(parents=True, exist_ok=True)
SETTINGS_DIR = Path(__file__).resolve().parent.parent / "settings"

@router.post("/upload_settings")
async def upload_settings(file: UploadFile = File(...)):
    filename = file.filename

    # Ensure filename is not None and ends in .json
    if not filename or not filename.endswith(".json"):
        return JSONResponse(status_code=400, content={"error": "Only .json files are allowed."})

    save_path = UPLOAD_DIR / filename  # Now type-safe

    try:
        with open(save_path, "wb") as f:
            content = await file.read()
            f.write(content)
    except Exception as e:
        return JSONResponse(status_code=500, content={"error": str(e)})

    return {"message": f"File uploaded successfully as {filename}"}

@router.get("/list_settings")
def list_settings():
    logger.debug(f"Looking for settings in: {SETTINGS_DIR}")
    files = [{"filename": f.name} for f in SETTINGS_DIR.glob("*.json") if f.is_file()]
    logger.debug(f"Found settings files: {files}")
    return {"files": files}

@router.get("/get_settings/{filename}")
async def get_settings(filename: str):
    file_path = SETTINGS_DIR / filename
    if not file_path.exists():
        raise HTTPException(status_code=404, detail="Settings file not found")
    with open(file_path, "r", encoding="utf-8") as f:
        data = json.load(f)
    return data


__all__ = ["router"]