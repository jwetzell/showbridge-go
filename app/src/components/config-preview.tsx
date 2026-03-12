import { useMemo } from "react"
import yaml from "js-yaml"
import type { Config } from "@/lib/config"
import { ScrollArea } from "@/components/ui/scroll-area"

type ConfigPreviewProps = {
  config: Config
}

export function ConfigPreview({ config }: ConfigPreviewProps) {
  const yamlString = useMemo(() => {
    // Build a clean config object for YAML output
    const output: Record<string, unknown> = {}

    if (config.modules.length > 0) {
      output.modules = config.modules.map((m) => {
        const mod: Record<string, unknown> = { id: m.id, type: m.type }
        if (m.params && Object.keys(m.params).length > 0) {
          mod.params = m.params
        }
        return mod
      })
    }

    if (config.routes.length > 0) {
      output.routes = config.routes.map((r) => {
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
      })
    }

    if (Object.keys(output).length === 0) {
      return "# Add modules and routes to generate config\nmodules: []\nroutes: []"
    }

    return yaml.dump(output, { lineWidth: -1, noRefs: true, sortKeys: false })
  }, [config])

  return (
    <div className="h-full flex flex-col">
      <div className="text-sm font-semibold pb-2">Config Preview</div>
      <ScrollArea className="flex-1 border rounded-md bg-card">
        <pre className="p-3 text-xs leading-relaxed whitespace-pre-wrap">{yamlString}</pre>
      </ScrollArea>
    </div>
  )
}
