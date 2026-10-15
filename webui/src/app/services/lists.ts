import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class ListsService {
  public processorListIds: string[] = [];

  registerProcessorList(path: string | undefined) {
    if (path === undefined) {
      return '';
    }

    const id = this.pathToId(path);
    if (!this.processorListIds.includes(id)) {
      this.processorListIds.push(id);
    }
    return id;
  }

  removeProcessorList(path: string | undefined) {
    if (path === undefined) {
      return;
    }

    const id = this.pathToId(path);
    this.processorListIds = this.processorListIds.filter((listId) => listId !== id);
  }

  pathToId(path: string) {
    return path.replaceAll('/', '.');
  }
}
