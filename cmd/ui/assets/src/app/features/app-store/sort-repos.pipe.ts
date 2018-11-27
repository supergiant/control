import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'sortRepos'
})
export class SortReposPipe implements PipeTransform {

  transform(items: any, args?: any): any {

    if(!items) return;

    const sortedItems = items.sort((itemA) => {
      if (itemA.name === 'supergiant') {
        return 1;
      }
      return -1;
    });

    return sortedItems;
  }

}
