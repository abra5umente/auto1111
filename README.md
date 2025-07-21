# Go Auto1111 Launcher & Web GUI

A unified toolkit for Stable Diffusion on your own hardware. Run it **headless** from the desktop or via a **browser‑friendly Web GUI**—both now share one code‑base and identical JSON‑driven settings.

---

## ✨ What’s New (July 2025)

| Area             | Added / Changed                                                                                                                                                             |
| ---------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Desktop CLI**  | • *Automatic* PNG naming in `/output` <br>• Robust error handling (non‑zero exit codes bubble up)                                                                           |
| **Web GUI**      | • Drag‑and‑drop **`.json` settings** upload <br>• Dropdown recall of uploaded presets <br>• Generates images with a single click (prompt + selected JSON)                   |
| **Config split** | • `.env` → **system‑level** stuff (paths, model folders, env BATs) <br>• `settings.json` → **generation parameters** (sampler, scheduler, width × height, steps, cfg scale) |

---

## 1. Directory Layout

```text
project/
│  README.md          ← you’re here
│  requirements.txt   ← Python deps for the Web GUI
│
├─ backend/           ← Go binary + FastAPI app
│   ├─ generator.exe  ← built Go binary
│   ├─ settings.json  ← symlink / copy chosen preset
│   ├─ output/        ← generated PNGs (auto‑timestamped)
│   └─ settings/      ← uploaded *.json presets
│
└─ auto1111_webapp/   ← simple HTML/CSS/JS front‑end
```

---

## 2. Installation & Build

### ⚙️ Prerequisites

* **Go 1.21+** (for the CLI generator)
* **Python 3.10+** (for FastAPI & AUTOMATIC1111)
* **Stable Diffusion WebUI** (AUTOMATIC1111) cloned locally
* *Windows users*: `ENVIRONMENT_BAT` & `AUTO1111_BAT` paths configured in `.env`

### 🐍 Python deps (Web GUI)

```bash
python -m venv venv
venv\Scripts\activate  # or `source venv/bin/activate`
pip install -r requirements.txt
```

`requirements.txt` generated today:

```text
fastapi
uvicorn[standard]
python-multipart
pydantic
```

### 🛠️ Build & Run — Desktop CLI

```bash
cd backend
# Build once
go build -o generator.exe .

# Single image
./generator.exe --prompt "a cat in space" --settings settings/sample.json
```

Images appear in `backend/output/`.

### 🌐 Build & Run — Web GUI

```bash
cd backend
# 1. build Go binary (same as above)
# 2. start FastAPI server (reload disabled for Go exe stability)
uvicorn main:app --host 0.0.0.0 --port 8000
```

Open [http://localhost:8000/app/](http://localhost:8000/app/) in your browser.

---

## 3. Using the Web GUI

1. **Upload** a `*.json` preset via **Settings → Upload**. The file is saved to `backend/settings/` and instantly appears in the dropdown.
2. **Select** the preset from the dropdown. Its parameters (sampler, width, etc.) are shown for confirmation.
3. Type your **prompt** (or leave blank to auto‑enhance) and press **Generate**.
4. The resulting PNG streams back and is also written to `backend/output/`.

> **Tip:** presets can be version‑controlled. Commit the `settings/` folder to Git for easy sharing.

---

## 4. Roadmap (July 2025)

| Feature                                    | Status                                                                      |
| ------------------------------------------ | --------------------------------------------------------------------------- |
| **Prompt History**                         | ✅ Planned – quick recall in GUI                                             |
| **Style Presets**                          | 🚧 Expanding – now driven by uploaded JSON; will ship with curated defaults |
| **Negative Prompt Generation**             | ✅ Keeps bad stuff out automatically                                         |
| **Model Selection**                        | ✅ Current dropdown; will obey `.env` / JSON overrides soon                  |
| **GUI Interface**                          | 🚧 MVP shipped; polishing UI & auth next                                    |
| **Prompt Templates (JSON import)**         | ✅ Core to new workflow                                                      |
| **Social Features** (share prompts / PNGs) | ✅ Still on the slate                                                        |
| **Containerisation** (Docker)              | ✅ Working on this                                                           |

---

## 5. Troubleshooting

| Symptom                            | Likely Cause                                 | Fix                                                                                                    |
| ---------------------------------- | -------------------------------------------- | ------------------------------------------------------------------------------------------------------ |
| `generator.exe failed: …`          | Invalid JSON (missing comma or wrong key)    | Validate with `jq . settings/your.json`                                                                |
| Image 512×512 despite bigger width | JSON keys wrong *or* values saved as strings | Keys are `IMAGE_WIDTH/IMAGE_HEIGHT` or `width/height` (caps or lower). Numbers **must not** be quoted. |
| `No .env file found` from AUTO1111 | Working dir incorrect for the child process  | Confirm `start_auto1111()` sets `cmd.Dir` to `backend/`                                                |
| PNG missing but no error           | Output path invalid                          | Ensure `output/` exists & that you have write permissions                                              |

---

## 6. FAQ

> **Q:** Do I need both `.env` *and* `settings.json`?
>
> **A:** Yes. Think of `.env` as *system settings* (where Python lives, which SD model to load, API keys). `settings.json` is *per‑image parameters* (sampler, size, steps). Separate files keep your secrets out of Git while letting you share artistic presets.

---

© 2025 MIT‑licensed.  Happy rendering! 🚀
