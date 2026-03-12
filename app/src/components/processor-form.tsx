import { useState, useEffect } from "react"
import type { SchemaType } from "@/lib/schema"
import type { ProcessorConfig } from "@/lib/config"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Label } from "@/components/ui/label"
import { ParamField } from "@/components/param-field"

type ProcessorFormProps = {
  processorTypes: SchemaType[]
  value: ProcessorConfig
  onChange: (processor: ProcessorConfig) => void
}

export function ProcessorForm({ processorTypes, value, onChange }: ProcessorFormProps) {
  const [selectedType, setSelectedType] = useState(value.type)
  const schema = processorTypes.find((t) => t.type === selectedType)

  useEffect(() => {
    setSelectedType(value.type)
  }, [value.type])

  const handleTypeChange = (newType: string) => {
    setSelectedType(newType)
    const newSchema = processorTypes.find((t) => t.type === newType)
    const defaultParams: Record<string, unknown> = {}
    if (newSchema?.params) {
      for (const param of newSchema.params) {
        if (param.default !== undefined) {
          defaultParams[param.name] = param.default
        }
      }
    }
    onChange({
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

  // For midi.message.create, handle oneOf variant params
  const hasOneOf = schema?.params?.some((p) => p.oneOf)
  const oneOfParam = schema?.params?.find((p) => p.oneOf)
  const selectedVariantType = oneOfParam ? (value.params?.type as string) : undefined
  const selectedVariant = oneOfParam?.oneOf?.find((v) => v.label === selectedVariantType)

  return (
    <div className="space-y-2">
      <div className="space-y-1">
        <Label className="text-xs">type *</Label>
        <Select value={selectedType} onValueChange={handleTypeChange}>
          <SelectTrigger className="h-9 text-xs">
            <SelectValue placeholder="Select processor type" />
          </SelectTrigger>
          <SelectContent>
            {processorTypes.map((t) => (
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

      {/* Render variant-specific params for oneOf (midi.message.create) */}
      {hasOneOf && selectedVariant && selectedVariant.params.map((param) => (
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
