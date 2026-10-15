import {
  CdkDrag,
  CdkDragDrop,
  CdkDragHandle,
  CdkDragPlaceholder,
  CdkDragPreview,
  CdkDropList,
  moveItemInArray,
} from '@angular/cdk/drag-drop';
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
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { cloneDeep } from 'lodash-es';
import { debounceTime, tap } from 'rxjs';
import { ProcessorConfig, RouteConfig } from '../../models/config';
import { ConfigService } from '../../services/config';
import { EventsService } from '../../services/events';
import { ListsService } from '../../services/lists';
import { SchemaService } from '../../services/schema';
import { ProcessorComponent } from '../processor/processor';

@Component({
  selector: 'app-route',
  imports: [
    MatFormFieldModule,
    MatIconModule,
    ReactiveFormsModule,
    MatMenuModule,
    JsonPipe,
    ProcessorComponent,
    MatTooltipModule,
    CdkDrag,
    CdkDragHandle,
    CdkDropList,
    CdkDragPreview,
    CdkDragPlaceholder,
  ],
  templateUrl: './route.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './route.css',
})
export class RouteComponent {
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
  route = input<RouteConfig>();
  moduleIds = input<string[]>([]);
  delete = output<void>();
  updated = output<RouteConfig>();
  moveProcessor = output<{
    fromRouteIndex: number;
    toRouteIndex: number;
    fromProcessorIndex: number;
    toProcessorIndex: number;
  }>();

  processors = computed(() => {
    const route = this.route();
    if (route) {
      return route.processors;
    }
    return [];
  });

  routeConfigErrors = computed(() => {
    const index = this.index();
    const routeErrors = this.configService.routeErrors();
    if (index !== undefined) {
      return routeErrors.filter((error) => error.index === index);
    }
    return [];
  });

  formGroup: FormGroup = new FormGroup({
    input: new FormControl('', [Validators.required]),
  });

  public schemaService = inject(SchemaService);
  private snackBar = inject(MatSnackBar);
  private eventsService = inject(EventsService);
  private configService = inject(ConfigService);
  public listsService = inject(ListsService);

  public processorListId = signal<string>('');

  indicatorColor = signal<string>('gray');

  ngOnDestroy(): void {
    this.listsService.removeProcessorList(this.path() + '/processors');
  }

  ngOnInit(): void {
    this.formGroup.patchValue({
      input: this.route()?.input,
    });

    this.formGroup.valueChanges.subscribe((value) => {
      const currentRoute = this.route();
      if (currentRoute === undefined) {
        console.error('route is undefined, not updating');
        return;
      }
      this.updated.emit({
        ...currentRoute,
        input: value.input,
      });
    });

    if (this.index() !== undefined) {
      this.eventsService
        .getRouteEventsForIndex(this.index()!)
        .pipe(
          tap((routeEvent) => {
            this.indicatorColor.set(routeEvent.error ? 'red' : 'greenyellow');
          }),
          debounceTime(100),
        )
        .subscribe((routeEvent) => {
          this.indicatorColor.set('gray');
        });
    }
    this.processorListId.set(this.listsService.registerProcessorList(this.path() + '/processors'));
  }

  isInError(): boolean {
    const path = this.path();
    if (path) {
      return this.schemaService.errorPaths.includes(path);
    }
    return false;
  }

  addProcessor(processorType: string) {
    const processorTemplate = this.schemaService.getSkeletonForProcessor(processorType);
    const currentRoute = this.route();
    if (currentRoute === undefined) {
      console.error('route is undefined, not updating');
      return;
    }
    const processors = currentRoute.processors || [];
    processors.push(processorTemplate);
    this.updated.emit({
      ...currentRoute,
      processors: cloneDeep(processors),
    });

    this.snackBar.open('Processor Added', 'Dismiss', {
      duration: 3000,
    });
  }

  processorUpdated(index: number, processor: ProcessorConfig | undefined) {
    if (processor === undefined) {
      console.error('processor is undefined, not updating');
      return;
    }
    const currentRoute = this.route();
    if (currentRoute === undefined) {
      console.error('route is undefined, not updating');
      return;
    }
    if (currentRoute.processors === undefined) {
      console.error('route processors is undefined, not updating');
      return;
    }
    currentRoute.processors[index].type = processor.type;
    if (processor.params !== undefined) {
      currentRoute.processors[index].params = processor.params;
    }
    this.updated.emit({
      ...currentRoute,
      processors: cloneDeep(currentRoute.processors),
    });
  }

  deleteProcessor(index: number) {
    const currentRoute = this.route();
    if (currentRoute === undefined) {
      console.error('route is undefined, not updating');
      return;
    }
    if (currentRoute.processors === undefined) {
      console.error('route processors is undefined, not updating');
      return;
    }
    currentRoute.processors.splice(index, 1);

    this.updated.emit({
      ...currentRoute,
      processors: cloneDeep(currentRoute.processors),
    });

    this.snackBar.open('Processor Removed', 'Dismiss', {
      duration: 3000,
    });
  }

  drop(event: CdkDragDrop<ProcessorConfig[] | undefined>) {
    if (event.previousContainer === event.container) {
      const currentRoute = this.route();
      if (currentRoute === undefined) {
        console.error('route is undefined, not updating');
        return;
      }
      const processors = currentRoute.processors;
      if (processors === undefined) {
        console.error('route processors is undefined, cannot move processor');
        return;
      }
      moveItemInArray(processors, event.previousIndex, event.currentIndex);
      this.updated.emit({
        ...currentRoute,
        processors: cloneDeep(processors),
      });
    } else {
      const fromRouteIndex = parseInt(event.previousContainer.id.split('.')[1], 10);
      const toRouteIndex = parseInt(event.container.id.split('.')[1], 10);
      const fromProcessorIndex = event.previousIndex;
      const toProcessorIndex = event.currentIndex;
      this.moveProcessor.emit({
        fromRouteIndex,
        toRouteIndex,
        fromProcessorIndex,
        toProcessorIndex,
      });
    }
  }
}
