import { SomeJSONSchema } from 'ajv/dist/types/json-schema';
import { ParamsFormInfo } from '../models/form';
import { FormGroup, FormControl } from '@angular/forms';
import { schemaToParamsFormInfo } from './params';

describe('params utils - schema to form', () => {
  const tests: { name: string; schema: SomeJSONSchema; formInfo: ParamsFormInfo }[] = [
    {
      name: 'basic string property',
      schema: {
        type: 'object',
        properties: {
          name: {
            title: 'Name',
            type: 'string',
          },
        },
        required: [],
      },
      formInfo: {
        formGroup: new FormGroup({
          name: new FormControl(''),
        }),
        paramsInfo: {
          name: {
            key: 'name',
            display: 'Name',
            type: 'string',
            hint: undefined,
            isConst: false,
            schema: {
              title: 'Name',
              type: 'string',
            },
            placeholder: '',
            default: undefined,
          },
        },
      },
    },
    {
      name: 'basic string property w/ default',
      schema: {
        type: 'object',
        properties: {
          name: {
            title: 'Name',
            type: 'string',
            default: 'John Doe',
          },
        },
        required: [],
      },
      formInfo: {
        formGroup: new FormGroup({
          name: new FormControl('John Doe'),
        }),
        paramsInfo: {
          name: {
            key: 'name',
            display: 'Name',
            type: 'string',
            hint: undefined,
            isConst: false,
            schema: {
              title: 'Name',
              type: 'string',
              default: 'John Doe',
            },
            placeholder: '',
            default: 'John Doe',
          },
        },
      },
    },
    {
      name: 'basic string property w/ description',
      schema: {
        type: 'object',
        properties: {
          name: {
            title: 'Name',
            description: 'Name of the thing',
            type: 'string',
          },
        },
        required: [],
      },
      formInfo: {
        formGroup: new FormGroup({
          name: new FormControl(''),
        }),
        paramsInfo: {
          name: {
            key: 'name',
            display: 'Name',
            type: 'string',
            hint: 'Name of the thing',
            isConst: false,
            schema: {
              title: 'Name',
              description: 'Name of the thing',
              type: 'string',
            },
            placeholder: '',
            default: undefined,
          },
        },
      },
    },
    {
      name: 'array of objects property',
      schema: {
        type: 'object',
        properties: {
          "args": {
            "description": "arguments",
            "items": {
              "additionalProperties": false,
              "properties": {
                "type": {
                  "description": "type of the argument",
                  "maxLength": 1,
                  "minLength": 1,
                  "title": "Type",
                  "type": "string"
                },
                "value": {
                  "description": "value of the argument",
                  "title": "Value",
                  "type": "string"
                }
              },
              "required": [
                "value",
                "type"
              ],
              "type": "object"
            },
            "title": "Arguments",
            "type": "array"
          },
        },
        required: [],
      },
      formInfo: {
        formGroup: new FormGroup({
          args: new FormControl(''),
        }),
        paramsInfo: {
          args: {
            key: 'args',
            display: 'Arguments',
            type: 'array',
            hint: 'arguments',
            isConst: false,
            schema: {
              title: 'Arguments',
              description: 'arguments',
              type: 'array',
              items: {
                type: 'object',
                additionalProperties: false,
                properties: {
                  type: {
                    description: 'type of the argument',
                    maxLength: 1,
                    minLength: 1,
                    title: 'Type',
                    type: 'string'
                  },
                  value: {
                    description: 'value of the argument',
                    title: 'Value',
                    type: 'string'
                  }
                },
                required: ['value', 'type']
              }
            },
            placeholder: '',
            default: undefined,
          },
        },
      },
    }
  ];

  tests.forEach((testCase) => {
    test(test.name, () => {
      const result = schemaToParamsFormInfo(testCase.schema);
      expect(result.paramsInfo).toEqual(testCase.formInfo.paramsInfo);
      expect(Object.keys(result.formGroup.controls)).toEqual(
        Object.keys(testCase.formInfo.formGroup.controls),
      );
      Object.entries(result.formGroup.controls).forEach(([key, control]) => {
        const expectedControl = testCase.formInfo.formGroup.controls[key];
        expect(expectedControl).toBeDefined();
        expect(control.value).toEqual(expectedControl?.value);
      });
    });
  });
});
