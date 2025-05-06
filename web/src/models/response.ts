export interface ValueResponse {
  status: number;
  value: unknown;
}

export interface ErrorResponse {
  status: number;
  error: string;
}

export interface SessionResponse {
  status: number;
  sessionId: string;
}

export interface FileUploadResponse {
  status: number;
  file: string;
}

export interface TaskStateResponse {
  status: number;
  task: string;
  taskState: string;
}
