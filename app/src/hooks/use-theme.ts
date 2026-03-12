import { useState, useEffect } from "react"

function getInitialDark(): boolean {
  const stored = localStorage.getItem("theme")
  if (stored === "dark") return true
  if (stored === "light") return false
  return window.matchMedia("(prefers-color-scheme: dark)").matches
}

export function useTheme() {
  const [dark, setDark] = useState(getInitialDark)

  useEffect(() => {
    document.documentElement.classList.toggle("dark", dark)
    localStorage.setItem("theme", dark ? "dark" : "light")
  }, [dark])

  const toggle = () => setDark((d) => !d)

  return { dark, toggle }
}
