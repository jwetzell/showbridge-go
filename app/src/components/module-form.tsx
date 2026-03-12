import { useState, useEffect } from "react"
import type { SchemaType } from "@/lib/schema"
import type { ModuleConfig } from "@/lib/config"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { ParamField } from "@/components/param-field"

type ModuleFormProps = {
  moduleTypes: SchemaType[]
  value: ModuleConfig
  onChange: (module: ModuleConfig) => void
}

export function ModuleForm({ moduleTypes, value, onChange }: ModuleFormProps) {
  const [selectedType, setSelectedType] = useState(value.type)
  const schema = moduleTypes.find((t) => t.type === selectedType)

  useEffect(() => {
    setSelectedType(value.type)
  }, [value.type])

  const handleTypeChange = (newType: string) => {
    setSelectedType(newType)
    const newSchema = moduleTypes.find((t) => t.type === newType)
    const defaultParams: Record<string, unknown> = {}
    if (newSchema?.params) {
      for (const param of newSchema.params) {
        if (param.default !== undefined) {
          defaultParams[param.name] = param.default
        }
      }
    }
    onChange({
      ...value,
      type: newType,
      params: newSchema?.params ? defaultParams : undefined,
    })
  }

  const handleParamChange = (name: string, paramValue: unknown) => {
    onChange({
      ...value,
      params: { ...(value.params || {}), [name]: paramValue },
    })
  }

  return (
    <div className="space-y-3">
      <div className="space-y-1">
        <Label className="text-xs">id *</Label>
        <Input
          className="h-9 text-xs"
          value={value.id}
          onChange={(e) => onChange({ ...value, id: e.target.value })}
          placeholder="Module ID"
        />
      </div>

      <div className="space-y-1">
        <Label className="text-xs">type *</Label>
        <Select value={selectedType} onValueChange={handleTypeChange}>
          <SelectTrigger className="h-9 text-xs">
            <SelectValue placeholder="Select module type" />
          </SelectTrigger>
          <SelectContent>
            {moduleTypes.map((t) => (
              <SelectItem key={t.type} value={t.type} className="text-xs">
                {t.type}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {schema?.params?.map((param) => (
        <ParamField
          key={param.name}
          schema={param}
          value={value.params?.[param.name]}
          onChange={(v) => handleParamChange(param.name, v)}
        />
      ))}
    </div>
  )
}
