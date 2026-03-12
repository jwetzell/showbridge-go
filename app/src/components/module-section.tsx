import type { SchemaType } from "@/lib/schema"
import type { ModuleConfig, ConfigAction } from "@/lib/config"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { ModuleForm } from "@/components/module-form"
import { Plus, Trash2 } from "lucide-react"

type ModuleSectionProps = {
  modules: ModuleConfig[]
  moduleTypes: SchemaType[]
  dispatch: React.Dispatch<ConfigAction>
}

export function ModuleSection({ modules, moduleTypes, dispatch }: ModuleSectionProps) {
  const handleAdd = () => {
    const firstType = moduleTypes[0]
    const defaultParams: Record<string, unknown> = {}
    if (firstType?.params) {
      for (const param of firstType.params) {
        if (param.default !== undefined) {
          defaultParams[param.name] = param.default
        }
      }
    }
    dispatch({
      type: "ADD_MODULE",
      module: {
        id: `module-${modules.length + 1}`,
        type: firstType?.type || "",
        params: firstType?.params ? defaultParams : undefined,
      },
    })
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h2 className="text-sm font-semibold">Modules</h2>
        <Badge variant="secondary" className="text-xs">
          {modules.length}
        </Badge>
      </div>

      {modules.map((module, index) => (
        <Card key={index}>
          <CardHeader className="pb-2 pt-3 px-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-xs font-medium">
                <Badge variant="outline" className="text-xs mr-2">
                  {module.type}
                </Badge>
                {module.id}
              </CardTitle>
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7"
                onClick={() => dispatch({ type: "REMOVE_MODULE", index })}
              >
                <Trash2 className="h-3 w-3" />
              </Button>
            </div>
          </CardHeader>
          <Separator />
          <CardContent className="pt-3 pb-3 px-3">
            <ModuleForm
              moduleTypes={moduleTypes}
              value={module}
              onChange={(updated) => dispatch({ type: "UPDATE_MODULE", index, module: updated })}
            />
          </CardContent>
        </Card>
      ))}

      <Button variant="outline" size="sm" className="w-full text-xs" onClick={handleAdd}>
        <Plus className="h-3 w-3 mr-1" />
        Add Module
      </Button>
    </div>
  )
}
