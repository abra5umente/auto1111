from fastapi import FastAPI, Request, Form
from fastapi.responses import FileResponse, JSONResponse
from fastapi.staticfiles import StaticFiles
from fastapi.middleware.cors import CORSMiddleware
from subprocess import run
import sys
from uuid import uuid4
from pathlib import Path
import json
from .routes.upload import router as upload_router
import logging

app = FastAPI()

BASE_DIR = Path(__file__).resolve().parent
FRONTEND_DIR = BASE_DIR.parent / "auto1111_webapp" / "app"
SETTINGS_DIR = BASE_DIR / "settings"
OUTPUT_DIR = BASE_DIR.parent / "output"
GENERATOR_EXE = BASE_DIR / "generator.exe"

# set up logging
logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

# Mount upload routes
app.include_router(upload_router)

# Mount static files for frontend
app.mount("/app", StaticFiles(directory=FRONTEND_DIR, html=True), name="frontend")

# Enable CORS (for dev use)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.post("/generate")
async def generate(
    prompt: str = Form(...),
    settings_file: str = Form(...),
):
    import shutil
    from subprocess import run
    from uuid import uuid4

    # 1️⃣  Make sure the chosen JSON exists
    src = SETTINGS_DIR / settings_file
    if not src.exists():
        return JSONResponse(status_code=400,
                            content={"error": f"{settings_file} not found"})

    # 2️⃣  Copy it to <backend>/settings.json for the Go binary
    generator_settings_path = GENERATOR_EXE.parent / "settings.json"
    shutil.copy(src, generator_settings_path)

    # 3️⃣  Build an output path (ensure directory exists)
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    output_file = OUTPUT_DIR / f"{uuid4()}.png"

    # 4️⃣  Launch the Go generator
    cmd = [
        str(GENERATOR_EXE),
        "--prompt",   prompt,
        "--output",   str(output_file),
        "--settings", str(generator_settings_path),
    ]
    logger.debug("Running command: %s", " ".join(cmd))

    result = run(cmd, stdout=sys.stdout, stderr=sys.stderr)
    if result.returncode != 0:
        logger.error("generator.exe failed: %s", result.stderr.decode())
        return JSONResponse(status_code=500,
                            content={"error": "Generation failed"})

    if not output_file.exists():
        logger.error("generator.exe reported success but file not found")
        return JSONResponse(status_code=500,
                            content={"error": "Output file missing"})

    # 5️⃣  Stream the PNG back to the client
    return FileResponse(output_file, media_type="image/png")
