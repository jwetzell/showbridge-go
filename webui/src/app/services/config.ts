import { HttpClient } from '@angular/common/http';
import { computed, effect, inject, Injectable, signal } from '@angular/core';
import { cloneDeep, isEqual } from 'lodash-es';
import {
  Config,
  ConfigError,
  ModuleConfig,
  ModuleError,
  RouteConfig,
  RouteError,
} from '../models/config';
import { EventsService } from './events';
import { SchemaService } from './schema';
import { SettingsService } from './settings';
@Injectable({
  providedIn: 'root',
})
export class ConfigService {
  pendingConfigIsValid = computed(() => {
    const config = this._currentlyShownConfig();
    if (config === undefined) {
      return false;
    }
    return this.schemaService.validate(this.schemaService.configSchemaId, config);
  });
  private runningConfig = signal<Config | undefined>(undefined);
  private _currentlyShownConfig = signal<Config | undefined>(undefined);
  readonly currentlyShownConfig = this._currentlyShownConfig.asReadonly();

  configIsDirty = computed(() => {
    const currentlyShown = this._currentlyShownConfig();
    const running = this.runningConfig();
    if (currentlyShown === undefined || running === undefined) {
      return false;
    }
    return !isEqual(currentlyShown, running);
  });

  moduleErrors = signal<ModuleError[]>([]);
  routeErrors = signal<RouteError[]>([]);

  private http = inject(HttpClient);
  private settingsService = inject(SettingsService);
  private eventsService = inject(EventsService);

  constructor(private schemaService: SchemaService) {
    effect(() => {
      console.log('config state changed', this._currentlyShownConfig());
    });

    effect(() => {
      switch (this.eventsService.status()) {
        case 'open':
          console.log('Websocket connection opened, reloading config');
          this.loadConfig();
          break;
      }
    });
  }

  loadConfig() {
    const configUrl = this.settingsService.configUrl();
    if (!configUrl) {
      console.error('Config URL is not set');
      this.setEmptyConfig();
      return;
    }
    this.http.get<Config>(configUrl.toString()).subscribe({
      next: (config) => {
        if (this.schemaService.validate(this.schemaService.configSchemaId, config)) {
          this.updateCurrentlyShownConfig(config);
          this.runningConfig.set(cloneDeep(config));
        } else {
          console.error('Config from server is invalid', config);
          this.setEmptyConfig();
        }
      },
      error: (err) => {
        console.error('Failed to load config from server', err);
        this.setEmptyConfig();
      },
    });
  }

  uploadConfig(config: Config) {
    if (this.schemaService.validate(this.schemaService.configSchemaId, config)) {
      this.updateCurrentlyShownConfig(config);
      const configUrl = this.settingsService.configUrl();
      if (!configUrl) {
        console.error('Config URL is not set, cannot upload config');
        return;
      }
      this.http.put<ConfigError>(configUrl.toString(), config).subscribe({
        next: () => {
          console.log('Config uploaded successfully');
          this.moduleErrors.set([]);
          this.routeErrors.set([]);
          this.runningConfig.set(cloneDeep(config));
        },
        error: (err) => {
          console.error('Problems occurred while uploading config', err);
          if (err.error) {
            const configError: ConfigError = err.error;
            this.moduleErrors.set(configError.moduleErrors ?? []);
            this.routeErrors.set(configError.routeErrors ?? []);
          }
          this.runningConfig.set(cloneDeep(config));
        },
      });
    } else {
      console.error('Uploaded config is invalid', config);
    }
  }

  setEmptyConfig() {
    console.log('Setting empty config');
    this.updateCurrentlyShownConfig({
      api: {
        enabled: true,
        port: 8080,
      },
      modules: [],
      routes: [],
    });
  }

  updateCurrentlyShownConfig(config: Config) {
    this._currentlyShownConfig.set(cloneDeep(config));
  }

  updateModules(modules: ModuleConfig[]) {
    var newConfig = cloneDeep(this._currentlyShownConfig());
    if (!newConfig) {
      console.error('No currently shown config to update modules on');
      return;
    }
    newConfig.modules = cloneDeep(modules);

    this.updateCurrentlyShownConfig(newConfig);
  }

  updateRoutes(routes: RouteConfig[]) {
    const currentConfig = this._currentlyShownConfig();
    if (!currentConfig) {
      console.error('No currently shown config to update routes on');
      return;
    }

    currentConfig.routes = cloneDeep(routes);

    this.updateCurrentlyShownConfig(currentConfig);
  }
}
