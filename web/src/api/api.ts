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

    const data = await response.json();
    if (data.session_id === undefined) {
      return "";
    }
    return data.session_id as string;
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

async function UploadMediaAsync(sessionId: string | null, formData: FormData) {
  if (!sessionId) {
    console.error("sessionId is null");
    return null;
  }
  try {
    const response = await fetch(`${API_URL}/medias`, {
      method: "POST",
      body: formData,
      headers: {
        "X-Session-Id": sessionId,
      },
    });

    if (response.status !== 200) {
      const responseBody = await response.json();
      console.error("failed to upload media : ", responseBody.error);
      return responseBody;
    }

    return await response.json();
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

export {
  GetSessionIdAsync,
  UploadMediaAsync,
  GetMediaTaskStatusAsync,
  GetUploadedMediaAsync,
};
