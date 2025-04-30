import React, { useState, useRef, useEffect } from "react";
import "./App.css";
import {
  GetMediaTaskStatusAsync,
  GetSubtitleAsync,
  GetUploadedMediaAsync,
  StartGenerateSubtitleTaskAsync,
  UploadMediaAsync,
} from "./api/api.ts";
import { RenewSessionIdAsync } from "./session/session.ts";
import { Box, Button, Text, Input } from "@chakra-ui/react";

function App() {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [taskStatus, setTaskStatus] = useState<string>("fetching");
  const [needRefreshTaskStatus, setNeedRefreshTaskStatus] =
    useState<boolean>(false);
  const [uploadedMedia, setUploadedMedia] = useState<string | null>("");
  const [waitingSubtitle, setWaitingSubtitle] = useState<boolean>(false);
  const [subtitle, setSubtitle] = useState<string | null>(null);

  const renewTaskStatus = () => {
    setNeedRefreshTaskStatus(true);
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setSelectedFile(file);
    }
  };

  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!selectedFile) {
      return;
    }
    const formData = new FormData();
    formData.append("file", selectedFile);
    const sessionId = sessionStorage.getItem("sessionId");
    if (!sessionId) {
      alert("No session id found.");
      return;
    }
    renewTaskStatus();
    UploadMediaAsync(sessionId, formData).then(async (response) => {
      if (response.error === "Session ID is expired") {
        setTaskStatus("Renewing session...");
        const newSession = await RenewSessionIdAsync();
        response = await UploadMediaAsync(newSession, formData);
      }
      if (response.error) {
        setTaskStatus("Error: " + response.error);
        alert(response.error);
        return;
      }
      alert("Uploaded");
      // refresh the uploaded media
      setUploadedMedia(null);
    });
  };

  const handleGenerateSubtitleClick = async () => {
    const response = await StartGenerateSubtitleTaskAsync();
    if (response === null || response.status !== 200) {
      if (response !== null) {
        const data = await response.json();
        if ("is_running" in data) {
          alert("Another task is running, please wait.");
        } else {
          alert("Error: " + data.error);
        }
        return;
      }
      alert(`Failed to start generating Subtitle`);
      return;
    }
    setWaitingSubtitle(true);
    renewTaskStatus();
  };

  const getSubtitle = async () => {
    const response = await GetSubtitleAsync();
    if (response === null || response.status !== 200) {
      let data = null;
      if (response) {
        data = await response.json();
      }
      alert(`Failed to get subtitle: ${data?.error}`);
      return;
    }
    const data = await response.json();
    if ("result" in data) {
      setSubtitle(data.result);
    }
  };

  const getUploadedMedia: () => Promise<string> = async () => {
    const response = await GetUploadedMediaAsync();
    if (response === null || response.status !== 200) {
      console.error("failed to get uploaded media: ", response?.status);
      return "";
    }
    const data = await response.json();
    if ("file_name" in data) {
      return data.file_name as string;
    }
    return "";
  };

  const getStatusAsync = async () => {
    const response = await GetMediaTaskStatusAsync();
    if (response === null || response.status !== 200) {
      console.error("failed to get task status : ", response?.status);
      return "Connection Failed";
    }
    const data = await response.json();
    if ("status" in data && data.status === "Running") {
      return data.task;
    }
    return "Idle";
  };

  useEffect(() => {
    const startTaskStatusTimer = () => {
      let tryCount = 0;
      const maxTryCount = 2;
      const intervalId = setInterval(async () => {
        const state = await getStatusAsync();
        setTaskStatus(state);
        if (state === "Idle" || state === "Connection Failed") {
          tryCount++;
          if (tryCount >= maxTryCount) {
            if (waitingSubtitle) {
              await getSubtitle();
              setWaitingSubtitle(false);
            }
            clearInterval(intervalId);
            setNeedRefreshTaskStatus(false);
          }
          return;
        }
        tryCount = 0;
      }, 1000);
    };
    if (needRefreshTaskStatus) {
      setTaskStatus("fetching task status...");
      startTaskStatusTimer();
    }
  }, [needRefreshTaskStatus, waitingSubtitle]);

  useEffect(() => {
    getUploadedMedia().then((response) => setUploadedMedia(response));
  }, [uploadedMedia]);

  return (
    <Box height="100vh" display="flex" flexDirection="column" maxHeight="100vh">
      <Text fontSize="lg" fontWeight="bold" textAlign="center" mt={4}>
        Task Status: {taskStatus}
      </Text>

      <Box flex={1} display="flex" flexDirection="row" overflow="auto">
        <Box
          display="flex"
          flexDirection="column"
          alignItems="flex-start"
          justifyContent="space-between"
          p={4}
          height="100%"
        >
          <form onSubmit={handleSubmit}>
            <Input
              type="file"
              ref={fileInputRef}
              onChange={handleFileChange}
              display="none"
            />

            <Button
              type="button"
              colorScheme="green"
              onClick={handleUploadClick}
              mb={2}
            >
              Select file
            </Button>

            {selectedFile && <Text mb={2}>{selectedFile.name}</Text>}

            <Button type="submit" colorScheme="blue" mb={2}>
              Upload file
            </Button>
          </form>

          <Button onClick={handleGenerateSubtitleClick} colorScheme="blue">
            Generate Subtitle
          </Button>
        </Box>

        <Box
          justifyContent="center"
          alignItems="self-start"
          height="100%"
          maxHeight="100%"
          overflow="auto"
          p={4}
        >
          <Text fontSize="lg" fontWeight="bold" whiteSpace="pre-wrap">
            {subtitle || "No subtitle generated yet."}
          </Text>
        </Box>
      </Box>
    </Box>
  );
}

export default App;
