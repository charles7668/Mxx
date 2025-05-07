# Mxx

🎬 **Mxx** is a user-friendly video subtitle generation tool that uses OpenAI's Whisper model to detect audio content and automatically generate subtitles.

> ⚠️ Currently supports **Linux systems only**

## ✨ Features

- 🗣️ **Automatic Speech Recognition (ASR)**  
  Transcribes spoken audio from videos into text using high-accuracy Whisper models.

- 🎞️ **Multi-Format Video Support**  
  Supports common video formats including MP4, MOV, MKV, and more.

- 📝 **Flexible Subtitle Formats**  
  Export subtitles as plain text or in the `ASS` subtitle format.

## 🚀 How to Use

> **Note:** Ensure that `ffmpeg` is installed on your system before using Mxx.

### 🌐 Web Interface

1. **Start the backend API:**

   ```bash
   mxx --web
   ```

2. **Create a `.env` file inside the `web/` directory with the following content:**

   ```text
   VITE_REACT_APP_API_URL=http://localhost:8080
   ```

3. **Start the frontend:**

   ```bash
   npm install
   npm run dev
   ```

This will start a local development server.
By default:

- The **frontend runs on port `5173`**
- The **backend API runs on port `8080`**

### 🧾 Command-Line Interface

To see available CLI options:

```bash
mxx --help
```

Use CLI commands to process video files directly from the terminal.

## 🔧 Build Guide

To build the project from source:

```bash
make build
```

- The backend executable will be generated as: `Mxx`
- The frontend static files will be built into the `dist/` folder
