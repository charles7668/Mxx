import { useState, useEffect } from "react";
import "./App.css";
import {
  GetMediaTaskStatusAsync,
  GetSubtitleAsync,
  StartGenerateSubtitleTaskAsync,
} from "./api/api.ts";
import { Box, Button, Text, Spinner, HStack } from "@chakra-ui/react";
import { TaskStatus } from "./models/task_status.ts";
import {
  ErrorResponse,
  TaskStateResponse,
  ValueResponse,
} from "./models/response.ts";
import SideMenu from "./components/SideMenu.tsx";
import { GenerateSubtitleRequest } from "./models/request.ts";

function App() {
  const [taskStatus, setTaskStatus] = useState<TaskStatus>({
    idle: true,
    message: "Idle",
  });
  const [needRefreshTaskStatus, setNeedRefreshTaskStatus] =
    useState<boolean>(true);
  const [waitingSubtitle, setWaitingSubtitle] = useState<boolean>(false);
  const [subtitle, setSubtitle] = useState<string | null>(null);

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
        />
        <Box height="100%" display="flex" flexDirection="column" flex="1">
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
            border="3px solid red"
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
