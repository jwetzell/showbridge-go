export type ParamSchema = {
  name: string
  type: "string" | "integer" | "number" | "boolean" | "array"
  required: boolean
  enum?: string[]
  default?: unknown
  minimum?: number
  maximum?: number
  minLength?: number
  maxLength?: number
  items?: { type: string }
  oneOf?: { label: string; params: ParamSchema[] }[]
}

export type SchemaType = {
  type: string
  title?: string
  params?: ParamSchema[]
  direction?: "input" | "output" | "both"
}

// "input" = route input only, "output" = route output only, "both" = either
const MODULE_DIRECTIONS: Record<string, "input" | "output" | "both"> = {
  "http.client": "both",
  "http.server": "input",
  "time.interval": "input",
  "time.timer": "both",
  "midi.input": "input",
  "midi.output": "output",
  "mqtt.client": "both",
  "nats.client": "both",
  "psn.client": "both",
  "serial.client": "both",
  "sip.call.server": "input",
  "sip.dtmf.server": "input",
  "net.tcp.client": "both",
  "net.tcp.server": "input",
  "net.udp.client": "output",
  "net.udp.multicast": "both",
  "net.udp.server": "input",
}

export function canBeInput(direction?: "input" | "output" | "both"): boolean {
  return direction !== "output"
}

export function canBeOutput(direction?: "input" | "output" | "both"): boolean {
  return direction !== "input"
}

type JsonSchema = {
  items?: {
    oneOf?: JsonSchemaEntry[]
  }
}

type JsonSchemaEntry = {
  title?: string
  properties?: Record<string, JsonSchemaProperty>
  required?: string[]
}

type JsonSchemaProperty = {
  const?: string
  type?: string
  enum?: string[]
  default?: unknown
  minimum?: number
  maximum?: number
  minLength?: number
  maxLength?: number
  properties?: Record<string, JsonSchemaProperty>
  required?: string[]
  items?: { type: string }
  oneOf?: JsonSchemaOneOfEntry[]
}

type JsonSchemaOneOfEntry = {
  type?: string
  properties?: Record<string, JsonSchemaProperty>
  required?: string[]
}

function extractParams(
  paramsSchema: JsonSchemaProperty,
  _parentRequired: boolean
): ParamSchema[] | undefined {
  if (!paramsSchema.properties) return undefined

  const requiredFields = paramsSchema.required || []
  const params: ParamSchema[] = []

  for (const [name, prop] of Object.entries(paramsSchema.properties)) {
    if (name === "type") continue // skip the type discriminator in nested oneOf params

    const param: ParamSchema = {
      name,
      type: (prop.type || "string") as ParamSchema["type"],
      required: requiredFields.includes(name),
      ...(prop.enum && { enum: prop.enum }),
      ...(prop.default !== undefined && { default: prop.default }),
      ...(prop.minimum !== undefined && { minimum: prop.minimum }),
      ...(prop.maximum !== undefined && { maximum: prop.maximum }),
      ...(prop.minLength !== undefined && { minLength: prop.minLength }),
      ...(prop.maxLength !== undefined && { maxLength: prop.maxLength }),
      ...(prop.items && { items: prop.items }),
    }

    params.push(param)
  }

  return params.length > 0 ? params : undefined
}

export function parseModuleSchema(schema: JsonSchema): SchemaType[] {
  const entries = schema.items?.oneOf || []
  return entries.map((entry) => {
    const typeConst = entry.properties?.type?.const
    const paramsSchema = entry.properties?.params as JsonSchemaProperty | undefined
    const required = entry.required || []

    const type = typeConst || "unknown"
    return {
      type,
      title: entry.title,
      params: paramsSchema ? extractParams(paramsSchema, required.includes("params")) : undefined,
      direction: MODULE_DIRECTIONS[type] || "both",
    }
  })
}

export function parseProcessorSchema(schema: JsonSchema): SchemaType[] {
  const entries = schema.items?.oneOf || []
  return entries.map((entry) => {
    const typeConst = entry.properties?.type?.const
    const paramsSchema = entry.properties?.params as JsonSchemaProperty | undefined
    const required = entry.required || []

    const result: SchemaType = {
      type: typeConst || "unknown",
    }

    if (paramsSchema) {
      // Handle oneOf (midi.message.create)
      if (paramsSchema.oneOf) {
        const oneOfParams = paramsSchema.oneOf.map((variant) => {
          const typeEnum = variant.properties?.type?.enum
          const label = typeEnum?.[0] || "unknown"
          const variantRequired = variant.required || []
          const params: ParamSchema[] = []

          if (variant.properties) {
            for (const [name, prop] of Object.entries(variant.properties)) {
              if (name === "type") continue
              params.push({
                name,
                type: (prop.type || "string") as ParamSchema["type"],
                required: variantRequired.includes(name),
                ...(prop.enum && { enum: prop.enum }),
                ...(prop.default !== undefined && { default: prop.default }),
              })
            }
          }

          return { label, params }
        })

        // Create a synthetic param for the midi type selector + mark the oneOf variants
        result.params = [
          {
            name: "type",
            type: "string",
            required: true,
            enum: oneOfParams.map((v) => v.label),
            oneOf: oneOfParams,
          },
        ]
      } else {
        result.params = extractParams(paramsSchema, required.includes("params"))
      }
    }

    return result
  })
}
