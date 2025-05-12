import { Box, Button, Input, Select, Spacer, Text } from "@chakra-ui/react";
import React, { useMemo, useRef, useState } from "react";
import { UploadMediaAsync } from "../api/api.ts";
import { ErrorResponse } from "../models/response.ts";
import { GenerateSubtitleRequest } from "../models/request.ts";

interface SideMenuProps {
  renewTaskStatus: () => void;
  onGenerateSubtitleClick: (request: GenerateSubtitleRequest) => void;
  onUploadedSuccess: () => void;
}

const SideMenu: React.FC<SideMenuProps> = ({
  renewTaskStatus,
  onGenerateSubtitleClick,
  onUploadedSuccess,
}) => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [whisperModel, setWhisperModel] = useState<string>("tiny");
  const [whisperLang, setWhisperLang] = useState<string>("auto");
  const supportWhisperModels = useMemo(() => {
    return [
      { value: "tiny", label: "Tiny (75 MiB)" },
      { value: "tiny.en", label: "Tiny EN (75 MiB)" },
      { value: "base", label: "Base (142 MiB)" },
      { value: "base.en", label: "Base EN (142 MiB)" },
      { value: "small", label: "Small (466 MiB)" },
      { value: "small.en", label: "Small EN (466 MiB)" },
      { value: "small.en-tdrz", label: "Small EN TDRZ (465 MiB)" },
      { value: "medium", label: "Medium (1.5 GiB)" },
      { value: "medium.en", label: "Medium EN (1.5 GiB)" },
      { value: "large-v1", label: "Large V1 (2.9 GiB)" },
      { value: "large-v2", label: "Large V2 (2.9 GiB)" },
      { value: "large-v2-q5_0", label: "Large V2 Q5_0 (1.1 GiB)" },
      { value: "large-v3", label: "Large V3 (2.9 GiB)" },
      { value: "large-v3-q5_0", label: "Large V3 Q5_0 (1.1 GiB)" },
      { value: "large-v3-turbo", label: "Large V3 Turbo (1.5 GiB)" },
      { value: "large-v3-turbo-q5_0", label: "Large V3 Turbo Q5_0 (547 MiB)" },
    ];
  }, []);

  const supportWhipserLangs = useMemo(() => {
    return [
      { value: "auto", label: "Auto" },
      { value: "en", label: "English" },
      { value: "zh", label: "Chinese" },
      { value: "ja", label: "Japanese" },
      { value: "ko", label: "Korean" },
      { value: "fr", label: "French" },
      { value: "de", label: "German" },
    ];
  }, []);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setSelectedFile(file);
    }
  };
  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };
  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!selectedFile) {
      return;
    }
    const formData = new FormData();
    formData.append("file", selectedFile);
    renewTaskStatus();
    const response = await UploadMediaAsync(formData);
    if (response === null) {
      alert(`Failed to upload file`);
      return;
    } else if (response.status !== 200) {
      const errResponse: ErrorResponse = await response.json();
      alert(`Failed to upload file : ` + errResponse.error);
      return;
    }
    alert("File uploaded successfully");
    onUploadedSuccess();
  };

  return (
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

        {selectedFile && (
          <Text
            maxW="100%"
            whiteSpace="normal"
            overflowWrap="break-word"
            wordBreak="break-word"
          >
            {selectedFile.name}
          </Text>
        )}

        <Button
          type="button"
          colorScheme="green"
          onClick={handleUploadClick}
          mb={2}
        >
          Select file
        </Button>

        <Button type="submit" colorScheme="blue" mb={2}>
          Upload file
        </Button>
      </form>

      <Spacer />

      <Text>Whisper Model</Text>
      <Select
        value={whisperModel}
        onChange={(e) => setWhisperModel(e.target.value)}
        mb={2}
      >
        {supportWhisperModels.map((model) => {
          return (
            <option key={model.value} value={model.value}>
              {model.value}
            </option>
          );
        })}
      </Select>
      <Text>Whisper Language</Text>
      <Select
        value={whisperLang}
        onChange={(e) => setWhisperLang(e.target.value)}
        mb={2}
      >
        {supportWhipserLangs.map((model) => {
          return (
            <option key={model.value} value={model.value}>
              {model.label}
            </option>
          );
        })}
      </Select>
      <Button
        onClick={() => {
          onGenerateSubtitleClick({
            Model: whisperModel,
            Language: whisperLang,
          });
        }}
        colorScheme="blue"
      >
        Generate Subtitle
      </Button>
    </Box>
  );
};

export default SideMenu;
