import { useState } from "react"
import { invoke } from "@tauri-apps/api/core"
import { save, open } from "@tauri-apps/plugin-dialog"
import type { Config, ConfigAction } from "@/lib/config"
import { useTheme } from "@/hooks/use-theme"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { FolderOpen, Save, Play, Square, Loader2, Sun, Moon } from "lucide-react"

type ToolbarProps = {
  config: Config
  dispatch: React.Dispatch<ConfigAction>
  running: boolean
  onStart: (configPath: string) => Promise<void>
  onStop: () => Promise<void>
}

function buildConfigOutput(config: Config): Record<string, unknown> {
  return {
    modules: config.modules.map((m) => {
      const mod: Record<string, unknown> = { id: m.id, type: m.type }
      if (m.params && Object.keys(m.params).length > 0) {
        mod.params = m.params
      }
      return mod
    }),
    routes: config.routes.map((r) => {
      const route: Record<string, unknown> = {
        input: r.input,
        output: r.output,
      }
      if (r.processors.length > 0) {
        route.processors = r.processors.map((p) => {
          const proc: Record<string, unknown> = { type: p.type }
          if (p.params && Object.keys(p.params).length > 0) {
            proc.params = p.params
          }
          return proc
        })
      }
      return route
    }),
  }
}

export function Toolbar({ config, dispatch, running, onStart, onStop }: ToolbarProps) {
  const [loading, setLoading] = useState(false)
  const [lastSavePath, setLastSavePath] = useState<string | null>(null)
  const { dark, toggle: toggleTheme } = useTheme()

  const handleLoad = async () => {
    const path = await open({
      filters: [{ name: "YAML", extensions: ["yaml", "yml"] }],
    })
    if (!path) return

    setLoading(true)
    try {
      const loaded = await invoke<Record<string, unknown>>("load_config", { path })
      const modules = (loaded.modules as Config["modules"]) || []
      const routes = ((loaded.routes as Config["routes"]) || []).map((r) => ({
        ...r,
        processors: r.processors || [],
      }))
      dispatch({ type: "LOAD_CONFIG", config: { modules, routes } })
      setLastSavePath(path)
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async () => {
    let path = lastSavePath
    if (!path) {
      const chosen = await save({
        filters: [{ name: "YAML", extensions: ["yaml", "yml"] }],
        defaultPath: "config.yaml",
      })
      if (!chosen) return
      path = chosen
    }

    setLoading(true)
    try {
      await invoke("save_config", { config: buildConfigOutput(config), path })
      setLastSavePath(path)
    } finally {
      setLoading(false)
    }
  }

  const handleRun = async () => {
    if (running) {
      await onStop()
      return
    }

    let configPath = lastSavePath
    if (!configPath) {
      const chosen = await save({
        filters: [{ name: "YAML", extensions: ["yaml", "yml"] }],
        defaultPath: "config.yaml",
      })
      if (!chosen) return
      configPath = chosen
    }

    setLoading(true)
    try {
      await invoke("save_config", { config: buildConfigOutput(config), path: configPath })
      setLastSavePath(configPath)
    } finally {
      setLoading(false)
    }

    await onStart(configPath)
  }

  return (
    <div className="flex items-center gap-2 px-3 py-2 border-b">
      <span className="text-sm font-bold mr-2">showbridge-go</span>
      <Separator orientation="vertical" className="h-5" />

      <Button
        variant="outline"
        size="sm"
        className="h-7 text-xs"
        onClick={handleLoad}
        disabled={loading}
      >
        <FolderOpen className="h-3 w-3 mr-1" />
        Load
      </Button>

      <Button
        variant="outline"
        size="sm"
        className="h-7 text-xs"
        onClick={handleSave}
        disabled={loading}
      >
        <Save className="h-3 w-3 mr-1" />
        Save
      </Button>

      <Separator orientation="vertical" className="h-5" />

      <Button
        variant={running ? "destructive" : "default"}
        size="sm"
        className="h-7 text-xs"
        onClick={handleRun}
        disabled={loading}
      >
        {loading ? (
          <Loader2 className="h-3 w-3 mr-1 animate-spin" />
        ) : running ? (
          <Square className="h-3 w-3 mr-1" />
        ) : (
          <Play className="h-3 w-3 mr-1" />
        )}
        {running ? "Stop" : "Run"}
      </Button>

      {lastSavePath && (
        <>
          <Separator orientation="vertical" className="h-5" />
          <span className="text-xs text-muted-foreground truncate max-w-xs">{lastSavePath}</span>
        </>
      )}

      <div className="ml-auto">
        <Button
          variant="ghost"
          size="sm"
          className="h-7 w-7 p-0"
          onClick={toggleTheme}
        >
          {dark ? <Sun className="h-3.5 w-3.5" /> : <Moon className="h-3.5 w-3.5" />}
        </Button>
      </div>
    </div>
  )
}
