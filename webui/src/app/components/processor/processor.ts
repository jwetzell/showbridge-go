import { CdkDragHandle, CdkDragPlaceholder, CdkDragPreview } from '@angular/cdk/drag-drop';
import { Component, computed, inject, input, output, ChangeDetectionStrategy } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { ProcessorConfig } from '../../models/config';
import { SchemaService } from '../../services/schema';
import { ParamsFormComponent } from '../params-form/params-form';

@Component({
  selector: 'app-processor',
  imports: [
    MatIconModule,
    ParamsFormComponent,
    ReactiveFormsModule,
    CdkDragHandle,
    CdkDragPreview,
    CdkDragPlaceholder,
  ],
  templateUrl: './processor.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './processor.css',
})
export class ProcessorComponent {
  path = input<string>('');
  processor = input<ProcessorConfig>();
  delete = output<void>();

  inDragList = input<boolean>(false);

  params = computed(() => this.processor()!.params);

  schema = computed(() => {
    return this.processor()?.type
      ? this.schemaService.getSchemaForProcessorType(this.processor()!.type)
      : undefined;
  });

  updated = output<ProcessorConfig>();

  hasParams = computed(() => {
    const schema = this.schema();
    return (
      schema !== undefined &&
      schema.properties !== undefined &&
      schema.properties.params !== undefined
    );
  });
  private schemaService = inject(SchemaService);

  paramsUpdated(params: any) {
    console.log('Params updated:', params);
    const currentProcessor = this.processor();
    if (currentProcessor !== undefined) {
      if (params !== undefined) {
        this.updated.emit({
          ...currentProcessor,
          params: params,
        });
      } else {
        this.updated.emit({
          ...currentProcessor,
        });
      }
    }
  }

  deleteMe() {
    this.delete.emit();
  }

  isInError(): boolean {
    const path = this.path();
    if (path) {
      return this.schemaService.errorPaths.includes(path);
    }
    return false;
  }
}
