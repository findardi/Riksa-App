import { refreshSession } from "./api";

type RefreshResult = Awaited<ReturnType<typeof refreshSession>>;

const inflight = new Map<string, Promise<RefreshResult>>();

export function refreshSinleFlight(refreshToken: string): Promise<RefreshResult> {
    const existing = inflight.get(refreshToken)
    if (existing) return existing;

    const pending = refreshSession(refreshToken).finally(() => {
        inflight.delete(refreshToken)
    })

    inflight.set(refreshToken, pending)
    return pending
}