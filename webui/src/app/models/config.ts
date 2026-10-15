export type Config = {
  api: ApiConfig;
  modules: ModuleConfig[];
  routes: RouteConfig[];
};

export type ConfigState = {
  config: Config;
  timestamp: number;
  isCurrent: boolean;
};

export type ApiConfig = {
  enabled: boolean;
  port: number;
};

export type ModuleConfig = {
  id?: string;
  type: string;
  params?: Record<string, any>;
};

export type RouteConfig = {
  id: string;
  input?: string;
  processors?: ProcessorConfig[];
};

export type ProcessorConfig = {
  id: string;
  type: string;
  params?: Record<string, any>;
};

export type ModuleError = {
  index: number;
  config: ModuleConfig;
  error: string;
};

export type RouteError = {
  index: number;
  config: RouteConfig;
  error: string;
};

export type ConfigError = {
  moduleErrors?: ModuleError[];
  routeErrors?: RouteError[];
};
