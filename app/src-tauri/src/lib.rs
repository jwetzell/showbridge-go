use serde::Serialize;
use std::io::{BufRead, BufReader};
use std::process::{Child, Command, Stdio};
use std::sync::Mutex;
use tauri::{Emitter, State};

struct AppState {
    child: Mutex<Option<Child>>,
    project_root: String,
    shell_path: String,
}

/// Get the user's full PATH by sourcing their login shell.
/// macOS GUI apps don't inherit the shell PATH, so we need to ask the shell for it.
fn get_shell_path() -> String {
    let shell = std::env::var("SHELL").unwrap_or_else(|_| "/bin/zsh".to_string());
    Command::new(&shell)
        .args(["-l", "-c", "echo $PATH"])
        .output()
        .ok()
        .and_then(|o| {
            if o.status.success() {
                String::from_utf8(o.stdout).ok().map(|s| s.trim().to_string())
            } else {
                None
            }
        })
        .unwrap_or_else(|| std::env::var("PATH").unwrap_or_default())
}

#[derive(Serialize)]
struct EnvironmentInfo {
    go_installed: bool,
    binary_exists: bool,
    go_version: Option<String>,
}

#[tauri::command]
fn check_environment(state: State<AppState>) -> EnvironmentInfo {
    let go_version = Command::new("go")
        .arg("version")
        .env("PATH", &state.shell_path)
        .output()
        .ok()
        .and_then(|o| {
            if o.status.success() {
                String::from_utf8(o.stdout).ok().map(|s| s.trim().to_string())
            } else {
                None
            }
        });

    let binary_path = std::path::Path::new(&state.project_root).join("showbridge");
    let binary_exists = binary_path.exists();

    EnvironmentInfo {
        go_installed: go_version.is_some(),
        binary_exists,
        go_version,
    }
}

#[derive(Serialize)]
struct BuildResult {
    success: bool,
    output: String,
}

#[tauri::command]
fn build_binary(state: State<AppState>) -> BuildResult {
    match Command::new("go")
        .args(["build", "-o", "showbridge", "./cmd/showbridge"])
        .current_dir(&state.project_root)
        .env("PATH", &state.shell_path)
        .output()
    {
        Ok(output) => {
            let stdout = String::from_utf8_lossy(&output.stdout).to_string();
            let stderr = String::from_utf8_lossy(&output.stderr).to_string();
            let combined = if stdout.is_empty() {
                stderr
            } else {
                format!("{}\n{}", stdout, stderr)
            };
            BuildResult {
                success: output.status.success(),
                output: combined.trim().to_string(),
            }
        }
        Err(e) => BuildResult {
            success: false,
            output: format!("Failed to run go build: {}", e),
        },
    }
}

#[derive(Serialize)]
struct Schemas {
    modules: serde_json::Value,
    processors: serde_json::Value,
    routes: serde_json::Value,
    config: serde_json::Value,
}

#[tauri::command]
fn read_schemas(state: State<AppState>) -> Result<Schemas, String> {
    let schema_dir = std::path::Path::new(&state.project_root).join("schema");

    let read_json = |name: &str| -> Result<serde_json::Value, String> {
        let path = schema_dir.join(name);
        let content =
            std::fs::read_to_string(&path).map_err(|e| format!("Failed to read {}: {}", name, e))?;
        serde_json::from_str(&content).map_err(|e| format!("Failed to parse {}: {}", name, e))
    };

    Ok(Schemas {
        modules: read_json("modules.schema.json")?,
        processors: read_json("processors.schema.json")?,
        routes: read_json("routes.schema.json")?,
        config: read_json("config.schema.json")?,
    })
}

#[tauri::command]
fn save_config(config: serde_json::Value, path: String) -> Result<(), String> {
    let yaml =
        serde_yaml::to_string(&config).map_err(|e| format!("Failed to serialize YAML: {}", e))?;
    std::fs::write(&path, yaml).map_err(|e| format!("Failed to write file: {}", e))
}

#[tauri::command]
fn load_config(path: String) -> Result<serde_json::Value, String> {
    let content =
        std::fs::read_to_string(&path).map_err(|e| format!("Failed to read file: {}", e))?;
    let value: serde_json::Value =
        serde_yaml::from_str(&content).map_err(|e| format!("Failed to parse YAML: {}", e))?;
    Ok(value)
}

#[tauri::command]
fn start_showbridge(
    config_path: String,
    state: State<AppState>,
    app_handle: tauri::AppHandle,
) -> Result<(), String> {
    let mut child_lock = state.child.lock().map_err(|e| e.to_string())?;

    if child_lock.is_some() {
        return Err("showbridge is already running".to_string());
    }

    let binary_path = std::path::Path::new(&state.project_root).join("showbridge");

    let mut child = Command::new(&binary_path)
        .args(["--config", &config_path, "--log-level", "info"])
        .current_dir(&state.project_root)
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
        .map_err(|e| format!("Failed to start showbridge: {}", e))?;

    let stdout = child.stdout.take();
    let stderr = child.stderr.take();

    *child_lock = Some(child);
    drop(child_lock);

    // Stream stdout
    if let Some(stdout) = stdout {
        let handle = app_handle.clone();
        std::thread::spawn(move || {
            let reader = BufReader::new(stdout);
            for line in reader.lines() {
                if let Ok(line) = line {
                    let _ = handle.emit("showbridge://log", &line);
                }
            }
        });
    }

    // Stream stderr
    if let Some(stderr) = stderr {
        let handle = app_handle.clone();
        std::thread::spawn(move || {
            let reader = BufReader::new(stderr);
            for line in reader.lines() {
                if let Ok(line) = line {
                    let _ = handle.emit("showbridge://log", &line);
                }
            }
        });
    }

    Ok(())
}

#[tauri::command]
fn stop_showbridge(state: State<AppState>) -> Result<(), String> {
    let mut child_lock = state.child.lock().map_err(|e| e.to_string())?;

    match child_lock.take() {
        Some(mut child) => {
            child.kill().map_err(|e| format!("Failed to kill process: {}", e))?;
            let _ = child.wait();
            Ok(())
        }
        None => Err("showbridge is not running".to_string()),
    }
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    // Resolve project root: go up from the app directory to find the Go project root
    // During `tauri dev`, CWD is app/src-tauri, so we need grandparent.
    // During bundled app, we look for the schema/ dir to confirm we found the right place.
    let project_root = std::env::current_dir()
        .ok()
        .and_then(|cwd| {
            // Try parent, grandparent, etc. looking for schema/ dir as anchor
            let mut dir = cwd.as_path();
            for _ in 0..5 {
                if dir.join("schema").is_dir() && dir.join("cmd").is_dir() {
                    return Some(dir.to_path_buf());
                }
                dir = dir.parent()?;
            }
            None
        })
        .unwrap_or_else(|| {
            // Fallback: assume we're in app/src-tauri
            let cwd = std::env::current_dir().unwrap_or_default();
            cwd.parent()
                .and_then(|p| p.parent())
                .map(|p| p.to_path_buf())
                .unwrap_or_else(|| std::path::PathBuf::from("../.."))
        })
        .to_string_lossy()
        .to_string();

    let shell_path = get_shell_path();

    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_dialog::init())
        .manage(AppState {
            child: Mutex::new(None),
            project_root,
            shell_path,
        })
        .setup(|app| {
            if cfg!(debug_assertions) {
                app.handle().plugin(
                    tauri_plugin_log::Builder::default()
                        .level(log::LevelFilter::Info)
                        .build(),
                )?;
            }
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            check_environment,
            build_binary,
            read_schemas,
            save_config,
            load_config,
            start_showbridge,
            stop_showbridge,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
