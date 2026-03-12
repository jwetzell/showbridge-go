import { useState, useEffect, useCallback, useRef } from "react"
import { invoke } from "@tauri-apps/api/core"
import { listen, type UnlistenFn } from "@tauri-apps/api/event"

export function useProcess() {
  const [running, setRunning] = useState(false)
  const [logs, setLogs] = useState<string[]>([])
  const unlistenRef = useRef<UnlistenFn | null>(null)

  useEffect(() => {
    return () => {
      if (unlistenRef.current) {
        unlistenRef.current()
      }
    }
  }, [])

  const start = useCallback(async (configPath: string) => {
    setLogs([])
    const unlisten = await listen<string>("showbridge://log", (event) => {
      setLogs((prev) => [...prev, event.payload])
    })
    unlistenRef.current = unlisten
    await invoke("start_showbridge", { configPath })
    setRunning(true)
  }, [])

  const stop = useCallback(async () => {
    await invoke("stop_showbridge")
    setRunning(false)
    if (unlistenRef.current) {
      unlistenRef.current()
      unlistenRef.current = null
    }
  }, [])

  const clearLogs = useCallback(() => {
    setLogs([])
  }, [])

  return { running, logs, start, stop, clearLogs }
}
