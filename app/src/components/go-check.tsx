import { useState, useEffect } from "react"
import { invoke } from "@tauri-apps/api/core"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import { AlertCircle, Hammer, Check, Loader2 } from "lucide-react"

type EnvironmentInfo = {
  go_installed: boolean
  binary_exists: boolean
  go_version: string | null
}

export function GoCheck() {
  const [env, setEnv] = useState<EnvironmentInfo | null>(null)
  const [building, setBuilding] = useState(false)
  const [buildOutput, setBuildOutput] = useState<string | null>(null)

  const checkEnv = async () => {
    const info = await invoke<EnvironmentInfo>("check_environment")
    setEnv(info)
  }

  useEffect(() => {
    checkEnv()
  }, [])

  const handleBuild = async () => {
    setBuilding(true)
    setBuildOutput(null)
    const result = await invoke<{ success: boolean; output: string }>("build_binary")
    setBuildOutput(result.output || (result.success ? "Build successful" : "Build failed"))
    if (result.success) {
      await checkEnv()
    }
    setBuilding(false)
  }

  if (!env) return null
  if (env.go_installed && env.binary_exists) return null

  return (
    <div className="space-y-2">
      {!env.go_installed && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle className="text-xs font-semibold">Go not found</AlertTitle>
          <AlertDescription className="text-xs">
            Go is required to build showbridge. Install it from{" "}
            <a
              href="https://go.dev/dl/"
              target="_blank"
              rel="noopener noreferrer"
              className="underline"
            >
              go.dev/dl
            </a>
          </AlertDescription>
        </Alert>
      )}

      {env.go_installed && !env.binary_exists && (
        <Alert>
          <Hammer className="h-4 w-4" />
          <AlertTitle className="text-xs font-semibold">Binary not built</AlertTitle>
          <AlertDescription className="text-xs space-y-2">
            <p>
              The showbridge binary needs to be built.
              {env.go_version && (
                <span className="text-muted-foreground ml-1">({env.go_version})</span>
              )}
            </p>
            <Button
              size="sm"
              className="h-7 text-xs"
              onClick={handleBuild}
              disabled={building}
            >
              {building ? (
                <Loader2 className="h-3 w-3 mr-1 animate-spin" />
              ) : (
                <Check className="h-3 w-3 mr-1" />
              )}
              {building ? "Building..." : "Build Binary"}
            </Button>
            {buildOutput && (
              <pre className="text-xs bg-muted p-2 rounded-md mt-2 whitespace-pre-wrap">
                {buildOutput}
              </pre>
            )}
          </AlertDescription>
        </Alert>
      )}
    </div>
  )
}
