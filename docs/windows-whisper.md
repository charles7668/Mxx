# Whisper GPU Support

## Windows

- Build [whisper.cpp](https://github.com/ggml-org/whisper.cpp) with GPU support enabled.
- Create `.a` files for MinGW by running the following commands for each DLL output by the Whisper build:

  ```bash
  gendef dll-name.dll
  dlltool -d dll-name.def -l dll-name.a -D dll-name.dll
  ```

- Copy the generated `.a` files to `whisper/lib/win`.
- When running the executable, ensure that all required `.dll` files are placed in the same directory as the `.exe`.
