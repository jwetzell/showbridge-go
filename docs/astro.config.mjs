import starlight from "@astrojs/starlight";
import { defineConfig } from "astro/config";

// https://astro.build/config
export default defineConfig({
  site: "https://docs.showbridge.io",
  integrations: [
    starlight({
      title: "showbridge",
      tableOfContents: false,
      favicon: "/favicon.ico",
      defaultLocale: 'en',
      editLink: {
        baseUrl: "https://github.com/jwetzell/showbridge-docs/edit/main/",
      },
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/jwetzell/showbridge-go",
        },
      ],
      sidebar: [
        {
          label: "showbridge",
          items: [{ autogenerate: { directory: "/showbridge" } }],
        },
        {
          label: "Run",
          items: [{ autogenerate: { directory: "run" } }],
        },
        {
          label: "Concepts",
          items: [{ autogenerate: { directory: "concepts" } }],
        },
        {
          label: "Modules",
          items: [
            {
              label: "HTTP",
              collapsed: true,
              items: [{ autogenerate: { directory: "modules/http" } }],
            },
            {
              label: "MIDI",
              collapsed: true,
              items: [{ autogenerate: { directory: "modules/midi" } }],
            },
            {
              label: "MQTT",
              collapsed: true,
              items: [{ autogenerate: { directory: "modules/mqtt" } }],
            },
            {
              label: "MySQL",
              collapsed: true,
              items: [{ autogenerate: { directory: "modules/mysql" } }],
            },
            {
              label: "NATS",
              collapsed: true,
              items: [{ autogenerate: { directory: "modules/nats" } }],
            },
            {
              label: "Network",
              collapsed: true,
              items: [
                {
                  label: "TCP",
                  collapsed: true,
                  items: [{ autogenerate: { directory: "modules/net/tcp" } }],
                },
                {
                  label: "UDP",
                  collapsed: true,
                  items: [{ autogenerate: { directory: "modules/net/udp" } }],
                },
              ],
            },
            {
              label: "PostgreSQL",
              collapsed: true,
              items: [{ autogenerate: { directory: "modules/postgres" } }],
            },
            {
              label: "PosiStageNet",
              collapsed: true,
              items: [{ autogenerate: { directory: "modules/psn" } }],
            },
            {
              label: "Redis",
              collapsed: true,
              items: [{ autogenerate: { directory: "modules/redis" } }],
            },
            {
              label: "Time",
              collapsed: true,
              items: [{ autogenerate: { directory: "modules/time" } }],
            },
            {
              label: "Serial",
              collapsed: true,
              items: [{ autogenerate: { directory: "modules/serial" } }],
            },
            {
              label: "SQLite",
              collapsed: true,
              items: [{ autogenerate: { directory: "modules/sqlite" } }],
            },
            {
              label: "WebSocket",
              collapsed: true,
              items: [{ autogenerate: { directory: "modules/websocket" } }],
            },
          ],
        },
        {
          label: "Processors",
          items: [
            {
              label: "ArtNet",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/artnet" } }],
            },
            {
              label: "Database",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/db" } }],
            },
            {
              label: "Debug",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/debug" } }],
            },
            {
              label: "Filter",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/filter" } }],
            },
            {
              label: "Float",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/float" } }],
            },
            {
              label: "FreeD",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/free-d" } }],
            },
            {
              label: "HTTP",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/http" } }],
            },
            {
              label: "Int",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/int" } }],
            },
            {
              label: "JSON",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/json" } }],
            },
            {
              label: "Key/Value",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/kv" } }],
            },
            {
              label: "MIDI",
              collapsed: true,
              items: [
                {
                  label: "Message",
                  collapsed: true,
                  items: [
                    { autogenerate: { directory: "processors/midi/message" } },
                  ],
                },
                {
                  label: "Note On",
                  collapsed: true,
                  items: [
                    { autogenerate: { directory: "processors/midi/note_on" } },
                  ],
                },
                {
                  label: "Note Off",
                  collapsed: true,
                  items: [
                    { autogenerate: { directory: "processors/midi/note_off" } },
                  ],
                },
                {
                  label: "Control Change",
                  collapsed: true,
                  items: [
                    {
                      autogenerate: {
                        directory: "processors/midi/control_change",
                      },
                    },
                  ],
                },
                {
                  label: "Program Change",
                  collapsed: true,
                  items: [
                    {
                      autogenerate: {
                        directory: "processors/midi/program_change",
                      },
                    },
                  ],
                },
              ],
            },
            {
              label: "Module",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/module" } }],
            },
            {
              label: "OS",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/os" } }],
            },
            {
              label: "OSC",
              collapsed: true,
              items: [
                {
                  label: "Message",
                  collapsed: true,
                  items: [
                    { autogenerate: { directory: "processors/osc/message" } },
                  ],
                },
              ],
            },
            {
              label: "Pub/Sub",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/pubsub" } }],
            },
            {
              label: "Route",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/route" } }],
            },
            {
              label: "Router",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/router" } }],
            },
            {
              label: "Script",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/script" } }],
            },
            {
              label: "String",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/string" } }],
            },
            {
              label: "Time",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/time" } }],
            },
            {
              label: "YAML",
              collapsed: true,
              items: [{ autogenerate: { directory: "processors/yaml" } }],
            },
          ],
        },
        {
          label: "Examples",
          items: [{ autogenerate: { directory: "examples" } }],
        },
      ],
    }),
  ],
});
