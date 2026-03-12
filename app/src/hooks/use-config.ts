import { useReducer } from "react"
import { configReducer, initialConfig } from "@/lib/config"

export function useConfig() {
  const [config, dispatch] = useReducer(configReducer, initialConfig)
  return { config, dispatch }
}
