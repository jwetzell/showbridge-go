import { CdkDrag, CdkDragDrop, CdkDropList, moveItemInArray } from '@angular/cdk/drag-drop';
import { Component, inject, Input, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { ParamInfo } from '../../models/form';
import { ListsService } from '../../services/lists';
import { SchemaService } from '../../services/schema';
import { parseStringToArray } from '../../utils/params';

@Component({
  selector: 'app-array-form',
  templateUrl: './array-form.html',
  styleUrl: './array-form.css',
  imports: [MatIconModule, FormsModule, CdkDrag, CdkDropList],
  changeDetection: ChangeDetectionStrategy.Eager,
  standalone: true,
})
export class ArrayFormComponent implements OnInit {
  @Input() paramFormControl?: any;
  @Input() paramInfo?: ParamInfo;

  minItems: number = 0;
  maxItems: number = Number.MAX_SAFE_INTEGER;

  arrayValue: any[] | undefined;

  private schemaService = inject(SchemaService);
  public listsService = inject(ListsService);

  constructor() {}

  ngOnInit(): void {
    if (this.paramFormControl && this.paramInfo?.schema) {
      if (!Array.isArray(this.paramFormControl.value)) {
        this.arrayValue = parseStringToArray(this.paramFormControl.value, this.paramInfo.schema);
      } else {
        this.arrayValue = this.paramFormControl.value;
      }

      if (this.paramInfo.schema.minItems) {
        this.minItems = parseInt(this.paramInfo?.schema.minItems);
      }

      if (this.paramInfo.schema.maxItems) {
        this.maxItems = parseInt(this.paramInfo?.schema.maxItems);
      }
    }
    this.ensureArrayMin();
  }

  dropItem(event: CdkDragDrop<any | undefined>) {
    if (this.arrayValue !== undefined) {
      moveItemInArray(this.arrayValue, event.previousIndex, event.currentIndex);
      this.valueUpdated();
    }
  }

  deleteItem(index: number) {
    this.arrayValue?.splice(index, 1);
    this.valueUpdated();
  }

  ensureArrayMin() {
    if (this.arrayValue === undefined) {
      this.arrayValue = [];
    }
    if (this.arrayValue.length < this.minItems) {
      for (let i = 0; i <= this.minItems - this.arrayValue.length; i += 1) {
        this.arrayValue.push(null);
      }
      this.valueUpdated();
    }
  }

  arrayIsMaxed() {
    if (this.arrayValue) {
      return this.arrayValue?.length >= this.maxItems;
    }
    return false;
  }

  addItem() {
    if (this.arrayValue === undefined) {
      this.arrayValue = [];
    }

    if (!this.arrayIsMaxed()) {
      this.arrayValue.push(null);
    }

    this.valueUpdated();
    // this.updated.emit(true);
  }

  valueUpdated() {
    if (this.arrayValue) {
      this.paramFormControl.setValue(
        this.schemaService.cleanArray(this.arrayValue, this.paramInfo?.schema?.items),
      );
    }
  }

  // NOTE(jwetzel): this is only needed for object item types
  updateItem(index: number, value: any) {
    if (this.arrayValue) {
      this.arrayValue[index] = value;
      this.valueUpdated();
    }
  }
}
