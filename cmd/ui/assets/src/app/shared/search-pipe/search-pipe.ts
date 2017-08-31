import { Pipe, PipeTransform } from '@angular/core';

// Tell Angular2 we're creating a Pipe with TypeScript decorators
@Pipe({
  name: 'Search'
})
export class Search implements PipeTransform {

  transform(items, args?) {
    if (args) {
      const search = args.toLowerCase();

      return items.filter(item => {
        return this.anyTrue(Object.getOwnPropertyNames(item)
          .map((key: string) => {
            if (item[key] && item[key].toString().toLowerCase().indexOf(search) >= 0) {
              return true;
            }
          }, {}));
      });
    }
    return items;
  }

  anyTrue(values) {
    for (const value of values) {
      if (value) {
        return true;
      }
    }
    return false;
  }
}
