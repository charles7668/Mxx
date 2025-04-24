import {GetSessionIdAsync} from "../api/api.ts";

async function RenewSessionIdAsync() {
    const newSession = await GetSessionIdAsync()
    if (!newSession) {
        console.warn("can not get sessionId");
    }
    sessionStorage.setItem("sessionId", newSession);
    return newSession;
}

export {RenewSessionIdAsync};