import { JsonPipe } from '@angular/common';
import {
  Component,
  effect,
  input,
  OnDestroy,
  output,
  ChangeDetectionStrategy,
} from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatTabsModule } from '@angular/material/tabs';
import { MatTooltipModule } from '@angular/material/tooltip';
import { SomeJSONSchema } from 'ajv/dist/types/json-schema';
import { cloneDeep, has, isEqual } from 'lodash-es';
import { Subscription } from 'rxjs';
import { ParamInfo, ParamsFormInfo } from '../../models/form';
import { cleanParams, schemaToParamsFormInfo } from '../../utils/params';
import { ArrayFormComponent } from '../array-form/array-form';
@Component({
  selector: 'app-params-form',
  templateUrl: './params-form.html',
  styleUrls: ['./params-form.css'],
  imports: [
    MatIconModule,
    MatTooltipModule,
    JsonPipe,
    MatFormFieldModule,
    MatInputModule,
    ReactiveFormsModule,
    ArrayFormComponent,
    MatTabsModule,
  ],
  changeDetection: ChangeDetectionStrategy.Eager,
  standalone: true,
})
export class ParamsFormComponent implements OnDestroy {
  data = input<any>();
  paramsSchema = input<SomeJSONSchema>();
  updated = output<any>();

  paramsFormInfo?: ParamsFormInfo;

  formGroupSubscription?: Subscription;
  constructor() {
    // TODO(jwetzell): do this in a bit more of a standard signal way
    effect(() => {
      const schema = this.paramsSchema();
      if (schema !== undefined) {
        if (schema.properties) {
          this.paramsFormInfo = schemaToParamsFormInfo(schema);
          if (this.formGroupSubscription === undefined) {
            this.formGroupSubscription = this.paramsFormInfo?.formGroup.valueChanges.subscribe(
              () => {
                this.formUpdated();
              },
            );
          }
        } else {
          console.error('params is not a singular object');
          console.error(schema);
        }
      }
    });

    effect(() => {
      const dataToPatch = cloneDeep(this.data());
      if (this.paramsSchema() && this.paramsFormInfo && dataToPatch) {
        if (this.paramsFormInfo?.formGroup !== undefined) {
          // NOTE(jwetzell): prepare data for form patching
          // const dataToPatch = cloneDeep(this.data);
          Object.entries(this.paramsFormInfo.paramsInfo).forEach(([paramKey, paramInfo]) => {
            if (has(dataToPatch, paramKey)) {
              switch (paramInfo.type) {
                case 'object':
                  dataToPatch[paramKey] = JSON.stringify(dataToPatch[paramKey]);
                  break;
                case 'array':
                  dataToPatch[paramKey] = dataToPatch[paramKey]
                    .map((item: any) => {
                      switch (typeof item) {
                        case 'object':
                          return JSON.stringify(item);
                        default:
                          return item;
                      }
                    })
                    .join(',');
                  break;
                default:
                  break;
              }
            }
          });

          this.paramsFormInfo.formGroup.patchValue(dataToPatch);
        }
      }
    });
  }

  ngOnDestroy(): void {
    if (this.formGroupSubscription) {
      this.formGroupSubscription.unsubscribe();
    }
  }

  formUpdated() {
    const paramsSchema = this.paramsSchema();
    if (paramsSchema) {
      const params = cleanParams(paramsSchema, this.paramsFormInfo?.formGroup.value);
      if (isEqual(params, this.data())) {
        // NOTE(jwetzell): no update
        return;
      }
      this.updated.emit(params);
    } else {
      console.error('params-form: no paramsSchema loaded');
    }
  }

  paramKeys() {
    if (this.paramsFormInfo) {
      return Object.keys(this.paramsFormInfo?.formGroup.controls);
    }
    return [];
  }

  getParamInfo(key: string): ParamInfo | undefined {
    return this.paramsFormInfo?.paramsInfo[key];
  }

  getParamValue(key: string) {
    const paramsSchema = this.paramsSchema();
    if (paramsSchema) {
      const params = cleanParams(paramsSchema, this.paramsFormInfo?.formGroup.value);
      return params[key];
    }
  }
}
