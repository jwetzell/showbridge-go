import { CdkDrag, CdkDragDrop, CdkDropList, moveItemInArray } from '@angular/cdk/drag-drop';
import { Component, inject, input, output, ChangeDetectionStrategy } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { cloneDeep } from 'lodash-es';
import { ModuleConfig } from '../../models/config';
import { ListsService } from '../../services/lists';
import { SchemaService } from '../../services/schema';
import { ModuleComponent } from '../module/module';

@Component({
  selector: 'app-module-list',
  imports: [
    MatMenuModule,
    MatIconModule,
    MatButtonModule,
    MatTooltipModule,
    ModuleComponent,
    CdkDrag,
    CdkDropList,
  ],
  templateUrl: './module-list.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './module-list.css',
})
export class ModuleListComponent {
  modules = input<ModuleConfig[]>();
  updated = output<ModuleConfig[]>();
  public schemaService = inject(SchemaService);
  public listService = inject(ListsService);

  private snackBar = inject(MatSnackBar);

  moduleUpdated(index: number, module: ModuleConfig | undefined) {
    if (module === undefined) {
      console.error('module is undefined, not updating');
      return;
    }
    const currentModules = this.modules();
    if (currentModules === undefined) {
      console.error('modules is undefined, cannot update');
      return;
    }
    currentModules[index] = cloneDeep(module);
    this.updated.emit(cloneDeep(currentModules));
  }

  addModule(moduleType: string) {
    const moduleTemplate = this.schemaService.getSkeletonForModule(moduleType);
    const currentModules = this.modules() || [];
    if (currentModules === undefined) {
      console.error('modules is undefined, cannot update');
      return;
    }
    currentModules.push(moduleTemplate);
    this.updated.emit(cloneDeep(currentModules));
    this.snackBar.open('Module Added', 'Dismiss', {
      duration: 3000,
    });
  }

  deleteModule(index: number) {
    const currentModules = this.modules();
    if (currentModules === undefined) {
      console.error('modules is undefined, cannot delete');
      return;
    }
    currentModules.splice(index, 1);
    this.updated.emit(cloneDeep(currentModules));
    this.snackBar.open('Module Removed', 'Dismiss', {
      duration: 3000,
    });
  }

  drop(event: CdkDragDrop<ModuleConfig[] | undefined>) {
    if (event.previousContainer === event.container) {
      const currentModules = this.modules();
      if (currentModules === undefined) {
        console.error('modules is undefined, not updating');
        return;
      }
      moveItemInArray(currentModules, event.previousIndex, event.currentIndex);
      this.updated.emit(cloneDeep(currentModules));
    }
  }
}
