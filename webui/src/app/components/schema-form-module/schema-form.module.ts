// 2. The Module File (user-management.module.ts)
import { NgModule } from '@angular/core';
import { CommonModule, JsonPipe } from '@angular/common';
import { ParamsFormComponent } from './params-form/params-form';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { CdkDrag, CdkDropList } from '@angular/cdk/drag-drop';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { ArrayFormComponent } from './array-form/array-form';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatTooltipModule } from '@angular/material/tooltip';

@NgModule({
    declarations: [
        ParamsFormComponent, // Declares the component belongs to this module
        ArrayFormComponent,

    ],
    imports: [
        CommonModule,
        MatIconModule,
        FormsModule,
        CdkDrag,
        CdkDropList,
        MatIconModule,
        MatTooltipModule,
        JsonPipe,
        MatFormFieldModule,
        MatInputModule,
        ReactiveFormsModule,
        MatTabsModule,
    ],
    exports: [
        ParamsFormComponent // Exports it so other modules can use it
    ]
})
export class SchemaFormModule { }