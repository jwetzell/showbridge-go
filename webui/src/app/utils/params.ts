import {
  AbstractControl,
  FormControl,
  FormGroup,
  ValidationErrors,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { SomeJSONSchema } from 'ajv/dist/types/json-schema';
import { ParamsFormInfo } from '../models/form';

export function schemaToParamsFormInfo(schema: SomeJSONSchema): ParamsFormInfo {
  const paramsFormInfo: ParamsFormInfo = {
    formGroup: new FormGroup({}),
    paramsInfo: {},
  };
  if (schema?.properties) {
    const paramKeys = Object.keys(schema.properties);
    Object.entries(schema.properties).forEach(([paramKey, paramSchema]: [string, any]) => {
      if (paramSchema.type) {
        switch (paramSchema.type) {
          case 'string':
          case 'number':
          case 'integer':
          case 'boolean':
          case 'array': // TODO(jwetzell): actually handle arrays
          case 'object': // TODO(jwetzell): actually handle objects
            let formDefault = '';
            const validators: ValidatorFn[] = [];

            // NOTE(jwetzell): check for a default value to set
            if (paramSchema.const) {
              formDefault = paramSchema.const;
            } else if (paramSchema.default) {
              formDefault = paramSchema.default;
            }

            // NOTE(jwetzell): add as many validators as we can
            if (paramSchema.minimum) {
              validators.push(Validators.min(paramSchema.minimum));
            }

            if (paramSchema.maximum) {
              validators.push(Validators.max(paramSchema.maximum));
            }

            if (paramSchema.pattern) {
              validators.push(Validators.pattern(new RegExp(paramSchema.pattern)));
            }

            if (schema.required) {
              if (schema.required.includes(paramKey)) {
                validators.push(Validators.required);
              }
            }

            paramsFormInfo.paramsInfo[paramKey] = {
              key: paramKey,
              display: paramSchema.title ? paramSchema.title : paramKey,
              type: paramSchema.type,
              hint: paramSchema.description ? paramSchema.description : undefined,
              isConst: !!paramSchema.const,
              schema: paramSchema,
              placeholder: '',
              default: paramSchema.default ? paramSchema.default : undefined,
            };

            if (paramSchema.examples && paramSchema.examples.length > 0) {
              paramsFormInfo.paramsInfo[paramKey].placeholder = paramSchema.examples[0];
            }

            if (paramSchema.type === 'object') {
              validators.push(objectValidator);
            }

            //TODO(jwetzell): figure out how to disable a control but not have to deal with undefined values on disabled controls
            paramsFormInfo.formGroup.addControl(paramKey, new FormControl(formDefault, validators));

            if (paramSchema.enum) {
              paramsFormInfo.paramsInfo[paramKey].options = paramSchema.enum;
            }

            break;
          default:
            console.error(
              `schema-service: unhandled param schema type for form group = ${paramSchema.type}`,
            );
            break;
        }
      } else {
        console.error('schema-service: param property without type');
      }
    });
  } else {
    console.error('trigger-form: params schema without properties');
    console.error(schema);
  }
  return paramsFormInfo;
}

export function cleanParams(paramsSchema: SomeJSONSchema, params: any): any {
  Object.keys(params).forEach((paramKey) => {
    if (paramsSchema.properties[paramKey]) {
      const paramSchema = paramsSchema.properties[paramKey];

      // delete null/undefined/empty params that aren't required
      if (params[paramKey] === undefined || params[paramKey] === null || params[paramKey] === '') {
        if (paramSchema.required) {
          if (!paramSchema.includes(paramKey)) {
            delete params[paramKey];
            return;
          }
        } else {
          delete params[paramKey];
          return;
        }
      }

      if (paramSchema.type) {
        switch (paramSchema.type) {
          case 'integer':
            var paramValue = parseInt(params[paramKey]);
            // NOTE(jwetzell): delete non-numbers
            if (Number.isNaN(paramValue)) {
              delete params[paramKey];
            } else {
              params[paramKey] = paramValue;
            }
            break;
          case 'number':
            var paramValue = parseFloat(params[paramKey]);
            // NOTE(jwetzell): delete non-numbers
            if (Number.isNaN(paramValue)) {
              delete params[paramKey];
            } else {
              params[paramKey] = paramValue;
            }
            break;
          case 'array':
            if (!Array.isArray(params[paramKey])) {
              const paramValue = params[paramKey];
              params[paramKey] = parseStringToArray(paramValue, paramSchema);
            }
            break;
          case 'object':
            try {
              params[paramKey] = JSON.parse(params[paramKey]);
            } catch (error) {
              console.error('object param is not JSON');
            }
            break;
          case 'string':
            // string is default so nothing needs to happen to clean it's value
            break;
          default:
            console.log(`schema-service: unhandled param schema type: ${paramSchema.type}`);
            break;
        }
      }
    }
  });
  return params;
}

function objectValidator(control: AbstractControl): ValidationErrors | null {
  try {
    JSON.parse(control.value);
    return null;
  } catch (error) {
    return { json: true };
  }
}

export function parseStringToArray(value: any, schema: SomeJSONSchema): any[] | undefined {
  if (!Array.isArray(value)) {
    const paramValue = value;
    if (paramValue === undefined || paramValue.trim().length === 0) {
      value = [];
    } else {
      if (schema.items?.type) {
        if (schema.items?.type === 'integer') {
          return paramValue
            .split(',')
            .map((part: string) => part.trim())
            .map((item: any) => parseInt(item));
        } else if (schema.items?.type === 'number') {
          return paramValue
            .split(',')
            .map((part: string) => part.trim())
            .map((item: any) => parseFloat(item));
        } else if (schema.items?.type === 'string') {
          return paramValue.split(',').map((part: string) => part.trim());
        } else if (schema.items?.type === 'object') {
          // TODO(jwetzell): this seems gross, not sure if this covers everything
          return JSON.parse(`[${paramValue}]`);
        } else {
          console.error(`schema-service: unhandled array schema type: ${schema.items?.type}`);
        }
      } else {
        // NOTE(jwetzell): default to comma-separated strings
        return paramValue.split(',').map((part: string) => part.trim());
      }
    }
  }
  return undefined;
}
