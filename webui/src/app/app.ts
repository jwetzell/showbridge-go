import { Component, computed, effect, inject, ChangeDetectionStrategy } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { ConfigPreviewComponent } from './components/config-preview/config-preview';
import { ModuleListComponent } from './components/module-list/module-list';
import { RouteListComponent } from './components/route-list/route-list';
import { Config, ModuleConfig, RouteConfig } from './models/config';
import { ConfigService } from './services/config';
import { EventsService } from './services/events';
import { SchemaService } from './services/schema';
import { SettingsService } from './services/settings';

@Component({
  selector: 'app-root',
  imports: [
    MatToolbarModule,
    MatProgressSpinnerModule,
    MatIconModule,
    MatMenuModule,
    MatButtonModule,
    MatTooltipModule,
    ModuleListComponent,
    RouteListComponent,
    ConfigPreviewComponent,
  ],
  templateUrl: './app.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './app.css',
})
export class App {
  config = computed<Config | undefined>(() => this.configService.currentlyShownConfig());

  modules = computed(() => this.config()?.modules ?? []);
  routes = computed(() => this.config()?.routes ?? []);

  moduleIds = computed(() => {
    const config = this.config();
    if (config !== undefined && config.modules !== undefined) {
      return config.modules.map((module) => module.id!) ?? [];
    }
    return [];
  });

  public schemaService = inject(SchemaService);
  public configService = inject(ConfigService);

  public settingsService = inject(SettingsService);
  public eventsService = inject(EventsService);
  private snackBar = inject(MatSnackBar);
  constructor() {
    effect(() => {
      if (this.schemaService.schemasLoaded()) {
        this.configService.loadConfig();
      }
    });
  }

  applyConfig() {
    const config = this.configService.currentlyShownConfig();
    if (config !== undefined) {
      this.configService.uploadConfig(config);
      this.snackBar.open('Config Uploaded', 'Dismiss', {
        duration: 3000,
      });
    }
  }

  loadConfig(config: Config) {
    this.configService.updateCurrentlyShownConfig(config);
  }

  modulesUpdated(modules: ModuleConfig[] | undefined) {
    if (modules === undefined) {
      console.error('modules is undefined, not updating');
      return;
    }
    this.configService.updateModules(modules);
  }

  routesUpdated(routes: RouteConfig[] | undefined) {
    if (routes === undefined) {
      console.error('routes is undefined, not updating');
      return;
    }
    this.configService.updateRoutes(routes);
  }
}
