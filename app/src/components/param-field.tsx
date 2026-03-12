import type { ParamSchema } from "@/lib/schema"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Button } from "@/components/ui/button"
import { Plus, X } from "lucide-react"

type ParamFieldProps = {
  schema: ParamSchema
  value: unknown
  onChange: (value: unknown) => void
}

function isEmptyRequired(schema: ParamSchema, value: unknown): boolean {
  if (!schema.required) return false
  if (value === undefined || value === null || value === "") return true
  if (Array.isArray(value) && value.length === 0) return true
  return false
}

function getStringError(schema: ParamSchema, value: unknown): string | undefined {
  if (typeof value !== "string" || value === "") return undefined
  if (schema.minLength !== undefined && schema.maxLength !== undefined && schema.minLength === schema.maxLength) {
    if (value.length !== schema.minLength) {
      return `Must be exactly ${schema.minLength} character${schema.minLength !== 1 ? "s" : ""}`
    }
  } else {
    if (schema.minLength !== undefined && value.length < schema.minLength) {
      return `Must be at least ${schema.minLength} character${schema.minLength !== 1 ? "s" : ""}`
    }
    if (schema.maxLength !== undefined && value.length > schema.maxLength) {
      return `Must be at most ${schema.maxLength} character${schema.maxLength !== 1 ? "s" : ""}`
    }
  }
  return undefined
}

export function ParamField({ schema, value, onChange }: ParamFieldProps) {
  const empty = isEmptyRequired(schema, value)

  // Handle oneOf (midi.message.create type selector)
  if (schema.oneOf) {
    const selectedType = (value as string) || ""
    const selectedVariant = schema.oneOf.find((v) => v.label === selectedType)

    return (
      <div className="space-y-2">
        <Label className="text-xs">{schema.name}{schema.required && " *"}</Label>
        <Select value={selectedType} onValueChange={onChange as (v: string) => void}>
          <SelectTrigger className={`h-9 text-xs ${empty ? "border-destructive" : ""}`}>
            <SelectValue placeholder={`Select ${schema.name}`} />
          </SelectTrigger>
          <SelectContent>
            {schema.oneOf.map((variant) => (
              <SelectItem key={variant.label} value={variant.label} className="text-xs">
                {variant.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {empty && <p className="text-destructive text-xs">Required</p>}
        {selectedVariant && selectedVariant.params.length > 0 && (
          <div className="text-xs text-muted-foreground pl-2">
            Fields: {selectedVariant.params.map((p) => p.name).join(", ")}
          </div>
        )}
      </div>
    )
  }

  // Enum → Select
  if (schema.enum) {
    return (
      <div className="space-y-1">
        <Label className="text-xs">{schema.name}{schema.required && " *"}</Label>
        <Select value={(value as string) || ""} onValueChange={onChange as (v: string) => void}>
          <SelectTrigger className={`h-9 text-xs ${empty ? "border-destructive" : ""}`}>
            <SelectValue placeholder={`Select ${schema.name}`} />
          </SelectTrigger>
          <SelectContent>
            {schema.enum.map((opt) => (
              <SelectItem key={opt} value={opt} className="text-xs">
                {opt}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {empty && <p className="text-destructive text-xs">Required</p>}
      </div>
    )
  }

  // Boolean → checkbox
  if (schema.type === "boolean") {
    return (
      <div className="flex items-center gap-2">
        <input
          type="checkbox"
          id={`param-${schema.name}`}
          checked={!!value}
          onChange={(e) => onChange(e.target.checked)}
          className="h-4 w-4"
        />
        <Label htmlFor={`param-${schema.name}`} className="text-xs">
          {schema.name}{schema.required && " *"}
        </Label>
      </div>
    )
  }

  // Array of strings → repeatable input
  if (schema.type === "array") {
    const items = (value as string[]) || []
    return (
      <div className="space-y-1">
        <Label className="text-xs">{schema.name}{schema.required && " *"}</Label>
        {items.map((item, i) => (
          <div key={i} className="flex gap-1">
            <Input
              className="h-9 text-xs"
              value={item}
              onChange={(e) => {
                const newItems = [...items]
                newItems[i] = e.target.value
                onChange(newItems)
              }}
              placeholder={`${schema.name} ${i + 1}`}
            />
            <Button
              variant="ghost"
              size="icon"
              className="h-9 w-9 shrink-0"
              onClick={() => {
                const newItems = items.filter((_, idx) => idx !== i)
                onChange(newItems)
              }}
            >
              <X className="h-3 w-3" />
            </Button>
          </div>
        ))}
        {empty && <p className="text-destructive text-xs">Required</p>}
        <Button
          variant="outline"
          size="sm"
          className="h-7 text-xs"
          onClick={() => onChange([...items, ""])}
        >
          <Plus className="h-3 w-3 mr-1" />
          Add
        </Button>
      </div>
    )
  }

  // Integer / Number → number input
  if (schema.type === "integer" || schema.type === "number") {
    return (
      <div className="space-y-1">
        <Label className="text-xs">{schema.name}{schema.required && " *"}</Label>
        <Input
          type="number"
          className={`h-9 text-xs ${empty ? "border-destructive" : ""}`}
          value={value !== undefined && value !== null ? String(value) : ""}
          min={schema.minimum}
          max={schema.maximum}
          onChange={(e) => {
            const v = e.target.value
            if (v === "") {
              onChange(undefined)
            } else {
              onChange(schema.type === "integer" ? parseInt(v, 10) : parseFloat(v))
            }
          }}
          placeholder={schema.default !== undefined ? String(schema.default) : undefined}
        />
        {empty && <p className="text-destructive text-xs">Required</p>}
      </div>
    )
  }

  // String → text input (default)
  const stringError = getStringError(schema, value)
  const hasError = empty || !!stringError

  return (
    <div className="space-y-1">
      <Label className="text-xs">{schema.name}{schema.required && " *"}</Label>
      <Input
        className={`h-9 text-xs ${hasError ? "border-destructive" : ""}`}
        value={(value as string) || ""}
        onChange={(e) => onChange(e.target.value)}
        placeholder={schema.default !== undefined ? String(schema.default) : undefined}
        minLength={schema.minLength}
        maxLength={schema.maxLength}
      />
      {empty && <p className="text-destructive text-xs">Required</p>}
      {stringError && <p className="text-destructive text-xs">{stringError}</p>}
    </div>
  )
}
