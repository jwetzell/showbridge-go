import { HttpClient } from '@angular/common/http';
import { computed, effect, inject, Injectable, signal } from '@angular/core';
import { AbstractControl, ValidationErrors } from '@angular/forms';
import { Ajv2020, JSONSchemaType } from 'ajv/dist/2020';
import { SomeJSONSchema } from 'ajv/dist/types/json-schema';
import { sortBy } from 'lodash-es';
import { Config, ModuleConfig, ProcessorConfig, RouteConfig } from '../models/config';
import { ObjectInfo } from '../models/form';
import { EventsService } from './events';
import { SettingsService } from './settings';

@Injectable({
  providedIn: 'root',
})
export class SchemaService {
  configSchema = signal<JSONSchemaType<Config> | undefined>(undefined);
  modulesSchema = signal<JSONSchemaType<ModuleConfig[]> | undefined>(undefined);
  routesSchema = signal<JSONSchemaType<RouteConfig[]> | undefined>(undefined);
  processorsSchema = signal<JSONSchemaType<ProcessorConfig[]> | undefined>(undefined);

  // TODO(jwetzell): populate from loaded schemas
  configSchemaId = 'https://showbridge.io/config.schema.json';
  modulesSchemaId = 'https://showbridge.io/modules.schema.json';
  routesSchemaId = 'https://showbridge.io/routes.schema.json';
  processorsSchemaId = 'https://showbridge.io/processors.schema.json';

  schemasLoaded = computed(() => {
    return (
      this.configSchema() !== undefined &&
      this.modulesSchema() !== undefined &&
      this.routesSchema() !== undefined &&
      this.processorsSchema() !== undefined
    );
  });

  ajv: Ajv2020 = new Ajv2020({ allErrors: true });

  processorTypes: ObjectInfo[] = [];
  moduleTypes: ObjectInfo[] = [];

  errorPaths: string[] = [];

  private settingsService = inject(SettingsService);
  private http = inject(HttpClient);
  private eventsService = inject(EventsService);

  constructor() {
    effect(() => {
      switch (this.eventsService.status()) {
        case 'open':
          if (this.settingsService.configSchemaUrl()) {
            this.loadConfigSchema();
          }
          if (this.settingsService.modulesSchemaUrl()) {
            this.loadModulesSchema();
          }
          if (this.settingsService.routesSchemaUrl()) {
            this.loadRoutesSchema();
          }
          if (this.settingsService.processorsSchemaUrl()) {
            this.loadProcessorsSchema();
          }
          break;
        case 'closed':
          this.resetAjv();
          this.configSchema.set(undefined);
          this.modulesSchema.set(undefined);
          this.routesSchema.set(undefined);
          this.processorsSchema.set(undefined);
          break;
      }
    });
  }

  resetAjv() {
    this.ajv.removeSchema(this.configSchema()!);
    this.ajv.removeSchema(this.modulesSchema()!);
    this.ajv.removeSchema(this.routesSchema()!);
    this.ajv.removeSchema(this.processorsSchema()!);
  }

  loadConfigSchema() {
    this.http
      .get<JSONSchemaType<Config>>(this.settingsService.configSchemaUrl().toString())
      .subscribe((schema) => {
        this.setConfigSchema(schema);
      });
  }

  loadModulesSchema() {
    this.http
      .get<JSONSchemaType<ModuleConfig[]>>(this.settingsService.modulesSchemaUrl().toString())
      .subscribe((schema) => {
        this.setModulesSchema(schema);
      });
  }

  loadRoutesSchema() {
    this.http
      .get<JSONSchemaType<RouteConfig[]>>(this.settingsService.routesSchemaUrl().toString())
      .subscribe((schema) => {
        this.setRoutesSchema(schema);
      });
  }

  loadProcessorsSchema() {
    this.http
      .get<JSONSchemaType<ProcessorConfig[]>>(this.settingsService.processorsSchemaUrl().toString())
      .subscribe((schema) => {
        this.setProcessorsSchema(schema);
      });
  }

  validate(schemaId: string, data: any): boolean {
    if (this.schemasLoaded() && this.ajv) {
      this.ajv.validate(schemaId, data);
      if (this.ajv.errors) {
        console.error('validation errors', this.ajv.errors);
        const errorPaths = new Set(
          this.ajv.errors.map((error) => {
            const errorPathMatch = error.instancePath.match(/(\/.*)\d+/);
            if (errorPathMatch) {
              return errorPathMatch[0]
                .substring(1)
                .replaceAll('handlers/', '')
                .replaceAll('/params', '');
            }
            return '';
          }),
        );
        this.errorPaths = Array.from(errorPaths);
        console.log(this.errorPaths);
        return false;
      } else {
        this.errorPaths = [];
      }
    }
    return true;
  }

  getSkeletonForRoute(): RouteConfig {
    const template: RouteConfig = {
      id: crypto.randomUUID(),
      processors: [],
    };
    return template;
  }

  getSkeletonForProcessor(processorType: string): ProcessorConfig {
    const template: ProcessorConfig = {
      id: crypto.randomUUID(),
      type: processorType,
    };
    const itemInfo = this.processorTypes.find((itemInfo) => itemInfo.type === processorType);
    if (itemInfo?.schema?.required && itemInfo.schema.required.includes('params')) {
      template.params = this.getSkeletonForParamsSchema(
        itemInfo.schema.properties.params,
      ) as Record<string, any>;
    }
    return template;
  }

  getSkeletonForModule(moduleType: string): ModuleConfig {
    const template: ModuleConfig = {
      type: moduleType,
    };
    const itemInfo = this.moduleTypes.find((itemInfo) => itemInfo.type === moduleType);
    if (itemInfo?.schema?.required && itemInfo.schema.required.includes('params')) {
      template.params = this.getSkeletonForParamsSchema(
        itemInfo.schema.properties.params,
      ) as Record<string, any>;
    }
    return template;
  }

  getSkeletonForParamsSchema(paramsSchema: SomeJSONSchema): Record<string, any> {
    const paramsTemplate: any = {};
    if (paramsSchema.properties) {
      Object.entries(paramsSchema.properties).forEach(([paramKey, paramSchema]) => {
        const schema = paramSchema as any;
        switch (schema.type) {
          case 'array':
            if (paramsSchema.required?.includes(paramKey)) {
              paramsTemplate[paramKey] = [];
            }
            break;
          default:
            console.error(`schema-service: unhandled param skeleton type = ${schema.type}`);
            break;
        }
      });
    }

    return paramsTemplate;
  }

  addSchemaToAjv(schema: SomeJSONSchema) {
    if (this.ajv.getSchema(schema.$id!)) {
      this.ajv.removeSchema(schema.$id!);
    }
    this.ajv.addSchema(schema);
  }

  setConfigSchema(schema: JSONSchemaType<Config>) {
    this.addSchemaToAjv(schema);
    this.configSchema.set(schema);
  }

  setModulesSchema(schema: JSONSchemaType<ModuleConfig[]>) {
    this.addSchemaToAjv(schema);
    this.modulesSchema.set(schema);
    this.moduleTypes = [];
    this.populateModuleTypes();
  }

  setRoutesSchema(schema: JSONSchemaType<RouteConfig[]>) {
    this.addSchemaToAjv(schema);
    this.routesSchema.set(schema);
  }

  setProcessorsSchema(schema: JSONSchemaType<ProcessorConfig[]>) {
    this.addSchemaToAjv(schema);
    this.processorsSchema.set(schema);
    this.processorTypes = [];
    this.populateProcessorTypes();
  }

  populateModuleTypes() {
    const modulesSchema = this.modulesSchema();
    if (modulesSchema !== undefined) {
      const definitions = modulesSchema.items?.oneOf;
      if (definitions) {
        this.moduleTypes = sortBy(
          Object.keys(definitions)
            .map((definitionKey) => definitions[definitionKey])
            .filter((definition) => definition.properties?.type?.const !== undefined)
            .map((definition) => {
              return {
                name: definition['title'] || definition.properties?.type?.const,
                type: definition.properties?.type?.const,
                schema: definition,
              };
            }),
          ['type'],
        );
      }
    }
  }

  populateProcessorTypes() {
    const processorsSchema = this.processorsSchema();
    if (processorsSchema !== undefined) {
      const definitions = processorsSchema.items?.oneOf;
      if (definitions) {
        this.processorTypes = sortBy(
          Object.keys(definitions)
            .map((definitionKey) => definitions[definitionKey])
            .filter((definition) => definition.properties?.type?.const !== undefined)
            .map((definition) => {
              return {
                name: definition['title'] || definition.properties?.type?.const,
                type: definition.properties?.type?.const,
                schema: definition,
              };
            }),
          ['type'],
        );
      }
    }
  }

  getSchemaForModuleType(protocolType: string) {
    const moduleType = this.moduleTypes.find((item) => item.type === protocolType);
    if (moduleType !== undefined) {
      return moduleType.schema;
    }
    return undefined;
  }

  getSchemaForProcessorType(protocolType: string) {
    const processorType = this.processorTypes.find((item) => item.type === protocolType);
    if (processorType !== undefined) {
      return processorType.schema;
    }
    return undefined;
  }

  matchParamsDataToSchema(data: any, schemas: SomeJSONSchema[]) {
    const matchingSchemaIndex = schemas.findIndex((schema) => {
      return this.ajv.validate(schema, data);
    });
    if (matchingSchemaIndex >= 0) {
      return matchingSchemaIndex;
    }
    return 0;
  }

  cleanArray(values: any[], itemSchema: SomeJSONSchema) {
    if (Array.isArray(values)) {
      switch (itemSchema.type) {
        case 'number':
          return values.map(Number.parseFloat);
        case 'integer':
          return values.map(Number.parseInt);
        case 'string':
          return values;
        case 'object':
          return values;
        default:
          console.error(`schema-service: unsupported array type ${itemSchema.type}`);
          return values;
      }
    }
    return [];
  }

  jsonValidator(validateSchema: SomeJSONSchema) {
    return (control: AbstractControl): ValidationErrors | null => {
      try {
        const jsonObj = JSON.parse(control.value);
        if (this.ajv.validate(validateSchema, jsonObj)) {
          return null;
        } else {
          return { json: this.ajv.errors };
        }
      } catch (error) {
        console.error(error);
        return { json: true };
      }
    };
  }
}
