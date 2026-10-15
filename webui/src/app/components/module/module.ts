import { CdkDragHandle, CdkDragPlaceholder, CdkDragPreview } from '@angular/cdk/drag-drop';
import { JsonPipe } from '@angular/common';
import {
  Component,
  computed,
  inject,
  input,
  output,
  signal,
  ChangeDetectionStrategy,
} from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule } from '@angular/material/tooltip';
import { cloneDeep } from 'lodash-es';
import { debounceTime, tap } from 'rxjs';
import { ModuleConfig } from '../../models/config';
import { ConfigService } from '../../services/config';
import { EventsService } from '../../services/events';
import { SchemaService } from '../../services/schema';
import { ParamsFormComponent } from '../params-form/params-form';

@Component({
  selector: 'app-module',
  imports: [
    MatFormFieldModule,
    MatIconModule,
    ParamsFormComponent,
    ReactiveFormsModule,
    MatMenuModule,
    MatTooltipModule,
    JsonPipe,
    CdkDragHandle,
    CdkDragPreview,
    CdkDragPlaceholder,
  ],
  templateUrl: './module.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './module.css',
})
export class ModuleComponent {
  path = input<string>('');
  index = computed(() => {
    const path = this.path();
    if (path) {
      const parts = path.split('/');
      const lastPart = parts[parts.length - 1];
      return parseInt(lastPart, 10);
    }
    return undefined;
  });

  inDragList = input<boolean>(false);
  module = input<ModuleConfig>();
  delete = output<void>();
  updated = output<ModuleConfig>();

  params = computed(() => this.module()!.params);
  id = computed(() => this.module()!.id);

  schema = computed(() => {
    return this.module()?.type
      ? this.schemaService.getSchemaForModuleType(this.module()!.type)
      : undefined;
  });

  moduleConfigErrors = computed(() => {
    const index = this.index();
    const moduleErrors = this.configService.moduleErrors();
    if (index !== undefined) {
      return moduleErrors.filter((error) => error.index === index);
    }
    return [];
  });

  formGroup: FormGroup = new FormGroup({
    id: new FormControl('', [Validators.required]),
    type: new FormControl(''),
  });

  inputIndicatorColor = signal<string>('gray');
  outputIndicatorColor = signal<string>('gray');

  private schemaService = inject(SchemaService);
  private eventsService = inject(EventsService);
  private configService = inject(ConfigService);
  ngOnInit(): void {
    this.formGroup.patchValue({
      id: this.module()?.id,
      type: this.module()?.type,
    });

    this.formGroup.valueChanges.subscribe((value) => {
      const currentModule = this.module();
      if (currentModule) {
        currentModule.id = value.id;
        currentModule.type = value.type;
        this.updated.emit(cloneDeep(currentModule));
      }
    });
    if (this.id() !== undefined) {
      this.eventsService
        .getInputEventsForSource(this.id()!)
        .pipe(
          tap((inputEvent) => {
            this.inputIndicatorColor.set(inputEvent.error ? 'red' : 'greenyellow');
          }),
          debounceTime(100),
        )
        .subscribe((inputEvent) => {
          this.inputIndicatorColor.set('gray');
        });
      this.eventsService
        .getOutputEventsForDestination(this.id()!)
        .pipe(
          tap((outputEvent) => {
            this.outputIndicatorColor.set(outputEvent.error ? 'red' : 'greenyellow');
          }),
          debounceTime(100),
        )
        .subscribe((outputEvent) => {
          this.outputIndicatorColor.set('gray');
        });
    }
  }

  paramsUpdated(params: any) {
    const currentModule = this.module();
    if (currentModule && currentModule.params !== undefined) {
      currentModule.params = params;
      this.updated.emit(cloneDeep(currentModule));
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
