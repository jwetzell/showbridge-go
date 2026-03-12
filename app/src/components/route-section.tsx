import type { SchemaType } from "@/lib/schema"
import { canBeInput, canBeOutput } from "@/lib/schema"
import type { ModuleConfig, RouteConfig, ConfigAction } from "@/lib/config"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { RouteForm } from "@/components/route-form"
import { Plus, Trash2 } from "lucide-react"

type ModuleOption = { id: string; type: string }

type RouteSectionProps = {
  routes: RouteConfig[]
  modules: ModuleConfig[]
  moduleTypes: SchemaType[]
  processorTypes: SchemaType[]
  dispatch: React.Dispatch<ConfigAction>
}

export function RouteSection({ routes, modules, moduleTypes, processorTypes, dispatch }: RouteSectionProps) {
  const typeDirectionMap = new Map(moduleTypes.map((t) => [t.type, t.direction]))

  const inputModules: ModuleOption[] = modules
    .filter((m) => canBeInput(typeDirectionMap.get(m.type)))
    .map((m) => ({ id: m.id, type: m.type }))

  const outputModules: ModuleOption[] = modules
    .filter((m) => canBeOutput(typeDirectionMap.get(m.type)))
    .map((m) => ({ id: m.id, type: m.type }))

  const handleAdd = () => {
    dispatch({
      type: "ADD_ROUTE",
      route: {
        input: inputModules[0]?.id || "",
        output: outputModules[0]?.id || "",
        processors: [],
      },
    })
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h2 className="text-sm font-semibold">Routes</h2>
        <Badge variant="secondary" className="text-xs">
          {routes.length}
        </Badge>
      </div>

      {routes.map((route, index) => (
        <Card key={index}>
          <CardHeader className="pb-2 pt-3 px-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-xs font-medium">
                {route.input} → {route.output}
                {route.processors.length > 0 && (
                  <Badge variant="secondary" className="text-xs ml-2">
                    {route.processors.length} processor{route.processors.length !== 1 && "s"}
                  </Badge>
                )}
              </CardTitle>
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7"
                onClick={() => dispatch({ type: "REMOVE_ROUTE", index })}
              >
                <Trash2 className="h-3 w-3" />
              </Button>
            </div>
          </CardHeader>
          <Separator />
          <CardContent className="pt-3 pb-3 px-3">
            <RouteForm
              routeIndex={index}
              value={route}
              inputModules={inputModules}
              outputModules={outputModules}
              processorTypes={processorTypes}
              dispatch={dispatch}
            />
          </CardContent>
        </Card>
      ))}

      <Button variant="outline" size="sm" className="w-full text-xs" onClick={handleAdd}>
        <Plus className="h-3 w-3 mr-1" />
        Add Route
      </Button>
    </div>
  )
}
