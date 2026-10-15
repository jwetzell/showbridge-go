import { A11yModule } from '@angular/cdk/a11y';
import {
  CdkDrag,
  CdkDragDrop,
  CdkDropList,
  moveItemInArray,
  transferArrayItem,
} from '@angular/cdk/drag-drop';
import { Component, inject, input, output, ChangeDetectionStrategy } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { cloneDeep } from 'lodash-es';
import { RouteConfig } from '../../models/config';
import { ListsService } from '../../services/lists';
import { SchemaService } from '../../services/schema';
import { RouteComponent } from '../route/route';

@Component({
  selector: 'app-route-list',
  imports: [
    RouteComponent,
    MatMenuModule,
    MatIconModule,
    MatButtonModule,
    MatTooltipModule,
    CdkDrag,
    CdkDropList,
    A11yModule,
  ],
  templateUrl: './route-list.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './route-list.css',
})
export class RouteListComponent {
  routes = input<RouteConfig[]>();
  moduleIds = input<string[]>();
  updated = output<RouteConfig[]>();
  public schemaService = inject(SchemaService);
  public listService = inject(ListsService);

  private snackBar = inject(MatSnackBar);

  deleteRoute(index: number) {
    const currentRoutes = this.routes();
    if (currentRoutes === undefined) {
      console.error('routes is undefined, cannot delete');
      return;
    }
    currentRoutes.splice(index, 1);
    this.updated.emit(cloneDeep(currentRoutes));
    this.snackBar.open('Route Removed', 'Dismiss', {
      duration: 3000,
    });
  }

  routeUpdated(index: number, route: RouteConfig | undefined) {
    if (route === undefined) {
      console.error('route is undefined, not updating');
      return;
    }
    const currentRoutes = this.routes();
    if (currentRoutes === undefined) {
      console.error('routes is undefined, cannot update');
      return;
    }
    currentRoutes[index] = cloneDeep(route);
    this.updated.emit(cloneDeep(currentRoutes));
  }

  addRoute() {
    const routeTemplate = this.schemaService.getSkeletonForRoute();
    const currentRoutes = this.routes() || [];
    if (currentRoutes === undefined) {
      console.error('routes is undefined, cannot update');
      return;
    }
    currentRoutes.push(routeTemplate);
    this.updated.emit(cloneDeep(currentRoutes));
    this.snackBar.open('Route Added', 'Dismiss', {
      duration: 3000,
    });
  }

  moveProcessorBetweenRoutes(event: {
    fromRouteIndex: number;
    toRouteIndex: number;
    fromProcessorIndex: number;
    toProcessorIndex: number;
  }) {
    if (event.fromRouteIndex === event.toRouteIndex) {
      console.error(
        'this should be handled by the route component, not moving processor between routes',
      );
      return;
    }
    const routes = this.routes();
    if (routes === undefined) {
      console.error('routes is undefined, cannot move processor');
      return;
    }
    const fromRoute = routes[event.fromRouteIndex];
    const toRoute = routes[event.toRouteIndex];
    if (fromRoute === undefined || toRoute === undefined) {
      console.error('fromRoute or toRoute is undefined, cannot move processor');
      return;
    }
    const fromProcessors = fromRoute.processors;
    const toProcessors = toRoute.processors;
    if (fromProcessors === undefined || toProcessors === undefined) {
      console.error('fromProcessors or toProcessors is undefined, cannot move processor');
      return;
    }
    transferArrayItem(
      fromProcessors,
      toProcessors,
      event.fromProcessorIndex,
      event.toProcessorIndex,
    );

    this.updated.emit(cloneDeep(routes));
  }

  drop(event: CdkDragDrop<RouteConfig[] | undefined>) {
    if (event.previousContainer === event.container) {
      const currentRoutes = this.routes();
      if (currentRoutes === undefined) {
        console.error('routes is undefined, not updating');
        return;
      }
      moveItemInArray(currentRoutes, event.previousIndex, event.currentIndex);
      this.updated.emit(cloneDeep(currentRoutes));
    }
  }
}
