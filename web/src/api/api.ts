import { SessionResponse } from "../models/response.ts";
import { GenerateSubtitleRequest } from "../models/request.ts";

const API_URL = import.meta.env.VITE_REACT_APP_API_URL || "";
console.log(`API URL: ${API_URL}`);

async function GetSessionIdAsync() {
  try {
    const response = await fetch(`${API_URL}/session`, {
      method: "GET",
    });

    if (response.status !== 200) {
      console.error("failed to get session id : ", response.status);
      return "";
    }

    const data: SessionResponse = await response.json();
    return data.sessionId;
  } catch (err) {
    console.error("Error fetching session ID:", err);
    return "";
  }
}

async function fetchWithAuth(requestUrl: string, request: RequestInit) {
  const storedSessionId = sessionStorage.getItem("sessionId") || "";
  const originHeaders = request.headers as Record<string, string>;
  let newHeaders = {
    ...originHeaders,
    "X-Session-Id": storedSessionId,
  };
  request.headers = newHeaders;
  const response = await fetch(requestUrl, request);
  if (response.status === 401) {
    const newSessionId = await GetSessionIdAsync();
    sessionStorage.setItem("sessionId", newSessionId);
    newHeaders = {
      ...originHeaders,
      "X-Session-Id": newSessionId,
    };
    request.headers = newHeaders;
    return await fetch(requestUrl, request);
  }
  return response;
}

async function GetMediaTaskStatusAsync() {
  try {
    return await fetchWithAuth(`${API_URL}/medias/task`, {
      method: "GET",
    });
  } catch (err) {
    console.error("Error fetching media task:", err);
    return null;
  }
}

async function UploadMediaAsync(formData: FormData) {
  try {
    return await fetchWithAuth(`${API_URL}/medias`, {
      method: "POST",
      body: formData,
    });
  } catch (err) {
    console.error("Error uploading media:", err);
    return null;
  }
}

async function GetUploadedMediaAsync() {
  try {
    return await fetchWithAuth(`${API_URL}/medias`, {
      method: "GET",
    });
  } catch (err) {
    console.error("Error fetching uploaded media:", err);
    return null;
  }
}

async function StartGenerateSubtitleTaskAsync(
  request: GenerateSubtitleRequest,
) {
  try {
    return await fetchWithAuth(`${API_URL}/medias/subtitles`, {
      method: "POST",
      body: JSON.stringify(request),
    });
  } catch (err) {
    console.error("Error starting subtitle task:", err);
    return null;
  }
}

async function GetSubtitleAsync() {
  try {
    return await fetchWithAuth(`${API_URL}/medias/subtitles`, {
      method: "GET",
    });
  } catch (err) {
    console.error("Error fetching subtitle:", err);
    return null;
  }
}

function GetPreviewMediaUrl() {
  const storedSessionId = sessionStorage.getItem("sessionId") || "";
  return `${API_URL}/video/${storedSessionId}/output.m3u8`;
}

export {
  GetSessionIdAsync,
  UploadMediaAsync,
  GetMediaTaskStatusAsync,
  GetUploadedMediaAsync,
  StartGenerateSubtitleTaskAsync,
  GetSubtitleAsync,
  GetPreviewMediaUrl,
};
