import { useState, useEffect } from "react";
import "./App.css";
import {
  GenerateSummary,
  GetASSFile,
  GetMediaTaskStatusAsync,
  GetPreviewMediaUrl,
  GetSubtitleAsync,
  StartGenerateSubtitleTaskAsync,
} from "./api/api.ts";
import { Box, Button, Text, Spinner, HStack } from "@chakra-ui/react";
import { TaskStatus } from "./models/task_status.ts";
import {
  ErrorResponse,
  SummaryResponse,
  TaskStateResponse,
  ValueResponse,
} from "./models/response.ts";
import SideMenu from "./components/SideMenu.tsx";
import {
  GenerateSubtitleRequest,
  GenerateSummaryRequest,
} from "./models/request.ts";
import VideoPlayer from "./components/VideoPlayer.tsx";

function App() {
  const [taskStatus, setTaskStatus] = useState<TaskStatus>({
    idle: true,
    message: "Idle",
  });
  const [needRefreshTaskStatus, setNeedRefreshTaskStatus] =
    useState<boolean>(true);
  const [waitingSubtitle, setWaitingSubtitle] = useState<boolean>(false);
  const [subtitle, setSubtitle] = useState<string | null>(null);
  const [videoUrl, setVideoUrl] = useState<string | null>(null);
  const [showPreviewVideo, setShowPreviewVideo] = useState<boolean>(false);
  const [showSummary, setShowSummary] = useState<boolean>(false);
  const [summary, setSummary] = useState<string | null>(null);

  const renewTaskStatus = () => {
    setNeedRefreshTaskStatus(true);
  };

  const handleGenerateSubtitleClick = async (
    request: GenerateSubtitleRequest,
  ) => {
    const response = await StartGenerateSubtitleTaskAsync(request);
    if (response === null) {
      alert("Failed to start generating Subtitle");
      return;
    } else if (response.status !== 200) {
      const data: ErrorResponse = await response.json();
      alert("Error: " + data.error);
      return;
    }
    setWaitingSubtitle(true);
    renewTaskStatus();
  };

  const getSubtitle = async () => {
    const response = await GetSubtitleAsync();
    if (response === null) {
      alert("Failed to get subtitle");
      return;
    } else if (response.status !== 200) {
      const data: ErrorResponse = await response.json();
      alert(`Failed to get subtitle: ${data.error}`);
      return;
    }
    const data: ValueResponse = await response!.json();
    setSubtitle(data.value as string);
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
    const data: TaskStateResponse = await response.json();
    if (data.taskState === "Running") {
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

  const showSummaryClickHandler = () => {
    setShowSummary(!showSummary);
  };

  const downloadASSClickHandler = async () => {
    const response = await GetASSFile();
    if (response === null) {
      alert("Failed to download ASS file");
      return;
    } else if (response.status !== 200) {
      const errResponse: ErrorResponse = await response.json();
      alert("Failed to download ASS file : " + errResponse.error);
      return;
    }

    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "subtitle.ass"; // Set the desired file name
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
  };

  const handleUploadSuccess = () => {
    setVideoUrl(GetPreviewMediaUrl());
  };

  const handleGenerateSummaryClick = async (
    request: GenerateSummaryRequest,
  ) => {
    renewTaskStatus();
    const response = await GenerateSummary(request);
    if (response === null) {
      alert("Failed to generate summary");
      return;
    }
    if (response.status !== 200) {
      const errResp: ErrorResponse = await response.json();
      alert(`Failed to generate summary: ${errResp.error}`);
      return;
    }
    const data: SummaryResponse = await response.json();
    setSummary(data.summary);
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
            clearInterval(intervalId);
            if (waitingSubtitle) {
              await getSubtitle();
              setWaitingSubtitle(false);
            }
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
        <SideMenu
          renewTaskStatus={renewTaskStatus}
          onGenerateSubtitleClick={handleGenerateSubtitleClick}
          onUploadedSuccess={handleUploadSuccess}
          onGenerateSummaryClick={handleGenerateSummaryClick}
        />
        <Box height="100%" display="flex" flexDirection="column" flex="1">
          {videoUrl && (
            <Button
              onClick={() => setShowPreviewVideo(!showPreviewVideo)}
              alignSelf="center"
              width="auto"
            >
              Toggle media preview
            </Button>
          )}

          {showPreviewVideo && (
            <Box mt={2} mb={2} ml="auto" mr="auto">
              <VideoPlayer videoUrl={videoUrl ?? ""}></VideoPlayer>
            </Box>
          )}

          <HStack>
            <Button onClick={showSummaryClickHandler}>
              {showSummary ? "Show Subtitle" : "Show summary"}
            </Button>
            <Button onClick={copyButtonHandler}>Copy</Button>
            <Button onClick={downloadASSClickHandler}>Download ASS</Button>
          </HStack>

          <Box
            justifyContent="center"
            alignItems="self-start"
            flex="1"
            maxHeight="100%"
            overflow="auto"
            p={4}
            border="3px solid red"
          >
            <Text fontSize="lg" fontWeight="bold" whiteSpace="pre-wrap">
              {showSummary && summary
                ? summary
                : subtitle || "No subtitle generated yet."}
            </Text>
          </Box>
        </Box>
      </Box>
    </Box>
  );
}

export default App;
