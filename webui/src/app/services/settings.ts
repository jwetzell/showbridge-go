import { computed, Injectable, signal } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class SettingsService {
  public configPath = signal('/api/v1/config');
  public configSchemaPath = signal('/schema/config.schema.json');
  public modulesSchemaPath = signal('/schema/modules.schema.json');
  public routesSchemaPath = signal('/schema/routes.schema.json');
  public processorsSchemaPath = signal('/schema/processors.schema.json');
  public wsPath = signal('/ws');
  public baseUrl = signal('http://localhost:8080');

  public configUrl = computed(() => {
    return new URL(this.configPath(), this.baseUrl());
  });

  public configSchemaUrl = computed(() => {
    return new URL(this.configSchemaPath(), this.baseUrl());
  });

  public modulesSchemaUrl = computed(() => {
    return new URL(this.modulesSchemaPath(), this.baseUrl());
  });

  public routesSchemaUrl = computed(() => {
    return new URL(this.routesSchemaPath(), this.baseUrl());
  });

  public processorsSchemaUrl = computed(() => {
    return new URL(this.processorsSchemaPath(), this.baseUrl());
  });

  public wsUrl = computed(() => {
    const url = new URL(this.baseUrl());
    url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:';
    url.pathname = this.wsPath();
    return url;
  });

  constructor() {
    const storedBaseUrl = localStorage.getItem('baseUrl');
    if (storedBaseUrl) {
      this.baseUrl.set(storedBaseUrl);
    }
  }

  updateBaseUrl(baseUrl: string) {
    this.baseUrl.set(baseUrl);
    localStorage.setItem('baseUrl', baseUrl);
  }
}
