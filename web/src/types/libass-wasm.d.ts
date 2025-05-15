declare module 'libass-wasm/dist/js/subtitles-octopus.js' {
  export default class SubtitlesOctopus {
    constructor(options: {
      video: HTMLVideoElement | null
      subContent: string
      fonts?: string[]
      workerUrl: string
      legacyWorkerUrl?: string
    }): void

    setTrack(content: string): void

    dispose(): void
  }
}
