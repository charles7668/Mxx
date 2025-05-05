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
import { Box, Button, Text, Input, Spinner, HStack } from "@chakra-ui/react";
import { TaskStatus } from "./models/task_status.ts";

function App() {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [taskStatus, setTaskStatus] = useState<TaskStatus>({
    idle: true,
    message: "Idle",
  });
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
        const newSession = await RenewSessionIdAsync();
        response = await UploadMediaAsync(newSession, formData);
      }
      if (response.error) {
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

  const getStatusAsync: () => Promise<TaskStatus> = async () => {
    const response = await GetMediaTaskStatusAsync();
    if (response === null || response.status !== 200) {
      console.error("failed to get task status : ", response?.status);
      return {
        idle: true,
        message: "Connection Failed",
      };
    }
    const data = await response.json();
    if ("status" in data && data.status === "Running") {
      return {
        idle: false,
        message: data.task,
      };
    }
    return {
      idle: true,
      message: "Idle",
    };
  };

  const copyButtonHandler = () => {
    if (subtitle) {
      navigator.clipboard.writeText(subtitle).then(() => {
        alert("Subtitle copied to clipboard");
      });
    } else {
      alert("No subtitle to copy");
    }
  };

  useEffect(() => {
    const startTaskStatusTimer = () => {
      let tryCount = 0;
      const maxTryCount = 2;
      const intervalId = setInterval(async () => {
        const state = await getStatusAsync();
        setTaskStatus(state);
        if (state.idle) {
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
      setTaskStatus({
        idle: false,
        message: "fetching task status...",
      });
      startTaskStatusTimer();
    }
  }, [needRefreshTaskStatus, waitingSubtitle]);

  useEffect(() => {
    getUploadedMedia().then((response) => setUploadedMedia(response));
  }, [uploadedMedia]);

  return (
    <Box height="100vh" display="flex" flexDirection="column">
      <Box
        display="flex"
        flexDirection="row"
        alignItems="center"
        justifyContent="center"
        alignContent="center"
        mt={4}
      >
        <Text fontSize="lg" fontWeight="bold" textAlign="center">
          Task Status: {taskStatus.message}
        </Text>
        {!taskStatus.idle && (
          <Box as="span" ml={2}>
            <Spinner size="sm" />
          </Box>
        )}
      </Box>

      <Box flex={1} display="flex" flexDirection="row" overflow="hidden">
        <Box
          display="flex"
          flexDirection="column"
          alignItems="flex-start"
          justifyContent="space-between"
          p={4}
          height="100%"
          maxHeight="100%"
          maxW="300px"
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

            {selectedFile && <Text whiteSpace="normal" mb={2}>{selectedFile.name}</Text>}

            <Button type="submit" colorScheme="blue" mb={2}>
              Upload file
            </Button>
          </form>

          <Button onClick={handleGenerateSubtitleClick} colorScheme="blue">
            Generate Subtitle
          </Button>
        </Box>
        <Box height="100%" display="flex" flexDirection="column">
          <HStack>
            <Button onClick={copyButtonHandler}>Copy</Button>
          </HStack>

          <Box
            justifyContent="center"
            alignItems="self-start"
            flex="1"
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
    </Box>
  );
}

export default App;
