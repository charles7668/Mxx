import Hls from "hls.js";
import React, { useEffect, useRef } from "react";

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
  }, [videoUrl]);

  return <video ref={videoRef} controls autoPlay muted width="600" />;
};

export default VideoPlayer;
