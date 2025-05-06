import { Box, Button, Input, Text } from "@chakra-ui/react";
import React, { useRef, useState } from "react";
import { UploadMediaAsync } from "../api/api.ts";
import { ErrorResponse } from "../models/response.ts";

interface SideMenuProps {
  renewTaskStatus: () => void;
  onGenerateSubtitleClick: () => void;
}

const SideMenu: React.FC<SideMenuProps> = ({
  renewTaskStatus,
  onGenerateSubtitleClick,
}) => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
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

        <Button
          type="button"
          colorScheme="green"
          onClick={handleUploadClick}
          mb={2}
        >
          Select file
        </Button>

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

        <Button type="submit" colorScheme="blue" mb={2}>
          Upload file
        </Button>
      </form>

      <Button onClick={onGenerateSubtitleClick} colorScheme="blue">
        Generate Subtitle
      </Button>
    </Box>
  );
};

export default SideMenu;
