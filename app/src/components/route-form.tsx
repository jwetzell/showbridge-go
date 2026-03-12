import type { SchemaType } from "@/lib/schema"
import type { RouteConfig, ProcessorConfig, ConfigAction } from "@/lib/config"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { ProcessorForm } from "@/components/processor-form"
import { Plus, Trash2, ChevronUp, ChevronDown } from "lucide-react"

type ModuleOption = { id: string; type: string }

type RouteFormProps = {
  routeIndex: number
  value: RouteConfig
  inputModules: ModuleOption[]
  outputModules: ModuleOption[]
  processorTypes: SchemaType[]
  dispatch: React.Dispatch<ConfigAction>
}

export function RouteForm({
  routeIndex,
  value,
  inputModules,
  outputModules,
  processorTypes,
  dispatch,
}: RouteFormProps) {
  const handleAddProcessor = () => {
    const firstType = processorTypes[0]
    const defaultParams: Record<string, unknown> = {}
    if (firstType?.params) {
      for (const param of firstType.params) {
        if (param.default !== undefined) {
          defaultParams[param.name] = param.default
        }
      }
    }
    dispatch({
      type: "ADD_PROCESSOR",
      routeIndex,
      processor: {
        type: firstType?.type || "",
        params: firstType?.params ? defaultParams : undefined,
      },
    })
  }

  return (
    <div className="space-y-3">
      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1">
          <Label className="text-xs">input *</Label>
          <Select
            value={value.input}
            onValueChange={(v) =>
              dispatch({ type: "UPDATE_ROUTE", index: routeIndex, route: { ...value, input: v } })
            }
          >
            <SelectTrigger className="h-9 text-xs">
              <SelectValue placeholder="Select input module" />
            </SelectTrigger>
            <SelectContent>
              {inputModules.map((m) => (
                <SelectItem key={m.id} value={m.id} className="text-xs">
                  {m.id} ({m.type})
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-1">
          <Label className="text-xs">output *</Label>
          <Select
            value={value.output}
            onValueChange={(v) =>
              dispatch({ type: "UPDATE_ROUTE", index: routeIndex, route: { ...value, output: v } })
            }
          >
            <SelectTrigger className="h-9 text-xs">
              <SelectValue placeholder="Select output module" />
            </SelectTrigger>
            <SelectContent>
              {outputModules.map((m) => (
                <SelectItem key={m.id} value={m.id} className="text-xs">
                  {m.id} ({m.type})
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      <Separator />

      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <Label className="text-xs font-medium">Processors</Label>
          <Badge variant="secondary" className="text-xs">
            {value.processors.length}
          </Badge>
        </div>

        {value.processors.map((processor, procIndex) => (
          <div key={procIndex} className="border rounded-md p-2 space-y-2">
            <div className="flex items-center justify-between">
              <Badge variant="outline" className="text-xs">
                #{procIndex + 1}
              </Badge>
              <div className="flex gap-1">
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-6 w-6"
                  disabled={procIndex === 0}
                  onClick={() =>
                    dispatch({
                      type: "MOVE_PROCESSOR",
                      routeIndex,
                      from: procIndex,
                      to: procIndex - 1,
                    })
                  }
                >
                  <ChevronUp className="h-3 w-3" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-6 w-6"
                  disabled={procIndex === value.processors.length - 1}
                  onClick={() =>
                    dispatch({
                      type: "MOVE_PROCESSOR",
                      routeIndex,
                      from: procIndex,
                      to: procIndex + 1,
                    })
                  }
                >
                  <ChevronDown className="h-3 w-3" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-6 w-6"
                  onClick={() =>
                    dispatch({ type: "REMOVE_PROCESSOR", routeIndex, processorIndex: procIndex })
                  }
                >
                  <Trash2 className="h-3 w-3" />
                </Button>
              </div>
            </div>
            <ProcessorForm
              processorTypes={processorTypes}
              value={processor}
              onChange={(updated: ProcessorConfig) =>
                dispatch({
                  type: "UPDATE_PROCESSOR",
                  routeIndex,
                  processorIndex: procIndex,
                  processor: updated,
                })
              }
            />
          </div>
        ))}

        <Button
          variant="outline"
          size="sm"
          className="w-full text-xs"
          onClick={handleAddProcessor}
        >
          <Plus className="h-3 w-3 mr-1" />
          Add Processor
        </Button>
      </div>
    </div>
  )
}
