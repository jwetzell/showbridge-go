export type ModuleConfig = {
  id: string
  type: string
  params?: Record<string, unknown>
}

export type ProcessorConfig = {
  type: string
  params?: Record<string, unknown>
}

export type RouteConfig = {
  input: string
  output: string
  processors: ProcessorConfig[]
}

export type Config = {
  modules: ModuleConfig[]
  routes: RouteConfig[]
}

export type ConfigAction =
  | { type: "ADD_MODULE"; module: ModuleConfig }
  | { type: "UPDATE_MODULE"; index: number; module: ModuleConfig }
  | { type: "REMOVE_MODULE"; index: number }
  | { type: "ADD_ROUTE"; route: RouteConfig }
  | { type: "UPDATE_ROUTE"; index: number; route: RouteConfig }
  | { type: "REMOVE_ROUTE"; index: number }
  | { type: "ADD_PROCESSOR"; routeIndex: number; processor: ProcessorConfig }
  | {
      type: "UPDATE_PROCESSOR"
      routeIndex: number
      processorIndex: number
      processor: ProcessorConfig
    }
  | { type: "REMOVE_PROCESSOR"; routeIndex: number; processorIndex: number }
  | { type: "MOVE_PROCESSOR"; routeIndex: number; from: number; to: number }
  | { type: "LOAD_CONFIG"; config: Config }

export function configReducer(state: Config, action: ConfigAction): Config {
  switch (action.type) {
    case "ADD_MODULE":
      return { ...state, modules: [...state.modules, action.module] }

    case "UPDATE_MODULE":
      return {
        ...state,
        modules: state.modules.map((m, i) => (i === action.index ? action.module : m)),
      }

    case "REMOVE_MODULE":
      return {
        ...state,
        modules: state.modules.filter((_, i) => i !== action.index),
      }

    case "ADD_ROUTE":
      return { ...state, routes: [...state.routes, action.route] }

    case "UPDATE_ROUTE":
      return {
        ...state,
        routes: state.routes.map((r, i) => (i === action.index ? action.route : r)),
      }

    case "REMOVE_ROUTE":
      return {
        ...state,
        routes: state.routes.filter((_, i) => i !== action.index),
      }

    case "ADD_PROCESSOR": {
      const routes = [...state.routes]
      const route = { ...routes[action.routeIndex] }
      route.processors = [...route.processors, action.processor]
      routes[action.routeIndex] = route
      return { ...state, routes }
    }

    case "UPDATE_PROCESSOR": {
      const routes = [...state.routes]
      const route = { ...routes[action.routeIndex] }
      route.processors = route.processors.map((p, i) =>
        i === action.processorIndex ? action.processor : p
      )
      routes[action.routeIndex] = route
      return { ...state, routes }
    }

    case "REMOVE_PROCESSOR": {
      const routes = [...state.routes]
      const route = { ...routes[action.routeIndex] }
      route.processors = route.processors.filter((_, i) => i !== action.processorIndex)
      routes[action.routeIndex] = route
      return { ...state, routes }
    }

    case "MOVE_PROCESSOR": {
      const routes = [...state.routes]
      const route = { ...routes[action.routeIndex] }
      const processors = [...route.processors]
      const [moved] = processors.splice(action.from, 1)
      processors.splice(action.to, 0, moved)
      route.processors = processors
      routes[action.routeIndex] = route
      return { ...state, routes }
    }

    case "LOAD_CONFIG":
      return action.config

    default:
      return state
  }
}

export const initialConfig: Config = {
  modules: [],
  routes: [],
}
