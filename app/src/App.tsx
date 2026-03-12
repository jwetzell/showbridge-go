import { useState, useEffect } from "react"
import { invoke } from "@tauri-apps/api/core"
import { parseModuleSchema, parseProcessorSchema } from "@/lib/schema"
import type { SchemaType } from "@/lib/schema"
import { useConfig } from "@/hooks/use-config"
import { useProcess } from "@/hooks/use-process"
import { Toolbar } from "@/components/toolbar"
import { GoCheck } from "@/components/go-check"
import { ModuleSection } from "@/components/module-section"
import { RouteSection } from "@/components/route-section"
import { ConfigPreview } from "@/components/config-preview"
import { LogPanel } from "@/components/log-panel"

type Schemas = {
  modules: unknown
  processors: unknown
  routes: unknown
  config: unknown
}

function App() {
  const { config, dispatch } = useConfig()
  const { running, logs, start, stop, clearLogs } = useProcess()
  const [moduleTypes, setModuleTypes] = useState<SchemaType[]>([])
  const [processorTypes, setProcessorTypes] = useState<SchemaType[]>([])

  useEffect(() => {
    invoke<Schemas>("read_schemas").then((schemas) => {
      setModuleTypes(parseModuleSchema(schemas.modules as Parameters<typeof parseModuleSchema>[0]))
      setProcessorTypes(
        parseProcessorSchema(schemas.processors as Parameters<typeof parseProcessorSchema>[0])
      )
    })
  }, [])

  return (
    <div className="h-screen flex flex-col">
      <Toolbar
        config={config}
        dispatch={dispatch}
        running={running}
        onStart={start}
        onStop={stop}
      />

      <div className="flex-1 overflow-hidden">
        <div className="h-full grid grid-cols-[1fr_1fr] gap-0">
          {/* Left: Config Editor */}
          <div className="overflow-y-auto p-4 space-y-6 border-r">
            <GoCheck />
            <ModuleSection modules={config.modules} moduleTypes={moduleTypes} dispatch={dispatch} />
            <RouteSection
              routes={config.routes}
              modules={config.modules}
              moduleTypes={moduleTypes}
              processorTypes={processorTypes}
              dispatch={dispatch}
            />
          </div>

          {/* Right: YAML Preview */}
          <div className="overflow-y-auto p-4">
            <ConfigPreview config={config} />
          </div>
        </div>
      </div>

      <LogPanel logs={logs} running={running} onClear={clearLogs} />
    </div>
  )
}

export default App
