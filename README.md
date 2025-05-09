# Mxx

🎬 **Mxx** is a user-friendly video subtitle generation tool that uses OpenAI's Whisper model to detect audio content and automatically generate subtitles.

## ✨ Features

- 🗣️ **Automatic Speech Recognition (ASR)**  
  Transcribes spoken audio from videos into text using high-accuracy Whisper models.

- 🎞️ **Multi-Format Video Support**  
  Supports common video formats including MP4, MOV, MKV, and more.

- 📝 **Flexible Subtitle Formats**  
  Export subtitles as plain text or in the `ASS` subtitle format.

## 🚀 How to Use

> **Note for Linux:** Make sure `ffmpeg` is installed and available in your system path.  
> **Note for Windows:** Place `ffmpeg.exe` in the same directory as `Mxx.exe`.

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

- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080`

### 🧾 Command-Line Interface

To view available CLI options:

```bash
mxx --help
```

Use CLI commands to process video files directly from the terminal.

## 🔧 Build Guide

> 💡 **Note:** It is recommended to compile the Whisper library yourself to enable more features such as GPU acceleration.  
> Otherwise, the default prebuilt version in the `whisper/lib` folder will be used (GPU is **not enabled** by default).

### 🪟 Notes for Windows

- The prebuilt Whisper library for Windows is located at:  
  `whisper/lib/win`

- If you encounter the error `0xc0000139` during runtime,  
  make sure to place the correct `libstdc++-6.dll` into the `whisper/lib/win` folder.

- You need to install **MinGW**, and ensure `make` and the required toolchains are available in your system `PATH`.

---

### 🛠️ Build Command

To build the project from source:

```bash
make build
```

- The backend executable will be generated as: `Mxx`
- The frontend static files will be output to the `dist/` folder
