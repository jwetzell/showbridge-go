export type RouterEvent<T, D> = {
  type: T;
  data?: D;
  error?: string;
};

export type RouteEventData = {
  index: number;
};

export type InputEventData = {
  source: string;
};

export type OutputEventData = {
  destination: string;
};
