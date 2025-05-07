import Hls from "hls.js";
import React, { useEffect, useRef } from "react";
import SubtitlesOctopus from "libass-wasm/dist/js/subtitles-octopus.js";
import { GetASSFile } from "../api/api.ts";

export interface VideoPlayerProps {
  videoUrl: string;
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({ videoUrl }) => {
  const videoRef = useRef<HTMLVideoElement | null>(null);

  useEffect(() => {
    if (Hls.isSupported()) {
      const hls = new Hls();
      hls.loadSource(videoUrl);
      hls.attachMedia(videoRef.current!);
    } else if (videoRef.current?.canPlayType("application/vnd.apple.mpegurl")) {
      videoRef.current.src = videoUrl;
    }
    GetASSFile().then(async (res) => {
      console.log(res);
      if (res === null || res.status !== 200) {
        return;
      }
      const assText = await res.text();
      console.log(assText);
      const options = {
        video: videoRef.current, // HTML5 video element
        subContent: assText,
        fonts: ["/fonts/default.woff2"], // Links to fonts (not required, default font already included in build)
        workerUrl: "/subtitles-octopus-worker.js", // Link to WebAssembly-based file "libassjs-worker.js"
        legacyWorkerUrl: "/subtitles-octopus-worker.js", // Link to non-WebAssembly worker
      };
      const instance = new SubtitlesOctopus(options);
      return () => {
        instance.dispose();
      };
    });
  }, [videoUrl]);

  return (
    <video
      ref={videoRef}
      controls
      muted
      width="600"
      style={{ maxHeight: "300px" }}
    />
  );
};

export default VideoPlayer;
