import { useRef, useEffect, useState } from "react"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { ChevronUp, ChevronDown, Trash2 } from "lucide-react"

type LogPanelProps = {
  logs: string[]
  running: boolean
  onClear: () => void
}

export function LogPanel({ logs, running, onClear }: LogPanelProps) {
  const [expanded, setExpanded] = useState(false)
  const scrollRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight
    }
  }, [logs])

  return (
    <div className="border-t">
      <div
        className="flex items-center justify-between px-3 py-1.5 cursor-pointer hover:bg-muted/50"
        onClick={() => setExpanded(!expanded)}
      >
        <div className="flex items-center gap-2">
          {expanded ? <ChevronDown className="h-3 w-3" /> : <ChevronUp className="h-3 w-3" />}
          <span className="text-xs font-medium">Logs</span>
          {running && (
            <Badge variant="default" className="text-[10px] h-4 px-1">
              running
            </Badge>
          )}
          {!running && logs.length > 0 && (
            <Badge variant="secondary" className="text-[10px] h-4 px-1">
              {logs.length} lines
            </Badge>
          )}
          {!expanded && logs.length > 0 && (
            <span className="text-xs text-muted-foreground truncate max-w-md">
              {logs[logs.length - 1]}
            </span>
          )}
        </div>
        {expanded && logs.length > 0 && (
          <Button
            variant="ghost"
            size="icon"
            className="h-6 w-6"
            onClick={(e) => {
              e.stopPropagation()
              onClear()
            }}
          >
            <Trash2 className="h-3 w-3" />
          </Button>
        )}
      </div>

      {expanded && (
        <ScrollArea className="h-48 border-t" ref={scrollRef}>
          <div className="p-2 space-y-0.5">
            {logs.length === 0 ? (
              <div className="text-xs text-muted-foreground py-4 text-center">
                No log output yet
              </div>
            ) : (
              logs.map((line, i) => (
                <div key={i} className="text-xs leading-relaxed font-mono">
                  {line}
                </div>
              ))
            )}
          </div>
        </ScrollArea>
      )}
    </div>
  )
}
